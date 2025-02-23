package redis

import (
  "context"
  "encoding/json"
  "errors"
  "github.com/go-redis/redis/v8"
  "log"
  "sync"
  "time"
)

type MultiLayerCache struct {
  client  *redis.Client
  localL1 sync.Map // L1 Cache (RAM local)
}

type ICache interface {
  Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
  Get(ctx context.Context, key string) (string, error)
  Delete(ctx context.Context, key string) error
}

func NewRedisCache(cfg *RedisConfig) *MultiLayerCache {
  rdb := redis.NewClient(&redis.Options{
    Addr:     cfg.Addr,
    Password: cfg.Password,
    Username: cfg.User,
  })

  // Check connection redis
  if err := rdb.Ping(context.Background()).Err(); err != nil {
    log.Printf("Warning: Unable to connect to Redis at %s: %v", cfg.Addr, err)
  } else {
    log.Println("Connected to Redis successfully!")
  }

  return &MultiLayerCache{client: rdb}
}

func (r *MultiLayerCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
  // Save into L1 Cache (Ram local)
  r.localL1.Store(key, value)

  // Save into L2 Cache (Redis)
  data, err := json.Marshal(value)
  if err != nil {
    log.Printf("❌ Failed to marshal cache data for key %s: %v", key, err)
    return err
  }

  err = r.client.Set(ctx, key, data, expiration).Err()
  if err != nil {
    log.Printf("⚠️ Redis set failed for key %s: %v", key, err)
    return err // Fast Fail nếu Redis lỗi
  }

  log.Printf("✅ Cache set success: %s (TTL: %v)", key, expiration)
  return nil
}

func (r *MultiLayerCache) Get(ctx context.Context, key string) (string, error) {
  if val, ok := r.localL1.Load(key); ok {
    log.Println("🔥 Cache Hit from L1 (RAM)")
    return val.(string), nil
  }

  // 2️⃣ if L1 not found , check L2 Cache (Redis)
  val, err := r.client.Get(ctx, key).Result()
  if errors.Is(err, redis.Nil) {
    return "", err
  } else if err != nil {
    // ❌ Fast Fail: Redis error
    log.Printf("⚠️ Redis error: %v (Fast Fail)", err)
    return "", err
  }
  log.Println("⚡ Cache Hit from L2 (Redis)")
  r.localL1.Store(key, val)
  return val, nil
}

func (r *MultiLayerCache) Delete(ctx context.Context, key string) error {
  r.localL1.Delete(key)               // Remove cache L1 cache
  return r.client.Del(ctx, key).Err() // Remove cache L2 cache
}
