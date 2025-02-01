package redis

import (
  "context"
  "github.com/go-redis/redis/v8"
  "time"
)

type RedisCache struct {
  client *redis.Client
}

type ICache interface {
  Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
  Get(ctx context.Context, key string) (string, error)
  Delete(ctx context.Context, key string) error
}

func NewRedisCache(cfg *RedisConfig) *RedisCache {
  rdb := redis.NewClient(&redis.Options{
    Addr:     cfg.Addr,
    Password: cfg.Password,
    Username: cfg.User,
  })

  return &RedisCache{client: rdb}
}

func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
  return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
  return r.client.Get(ctx, key).Result()
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
  return r.client.Del(ctx, key).Err()
}
