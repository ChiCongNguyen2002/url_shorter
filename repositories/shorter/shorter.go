package shorter

import (
  "context"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  "time"
  mongodb2 "url-shortener/database/mongodb"
)

type shorterRepository struct {
  *mongodb2.Repository[URLs]
}

type IShorterRepository interface {
  SaveURL(ctx context.Context, shortKey, longURL string, expireIn time.Duration) error
  GetURL(ctx context.Context, shortKey string) (*URLs, error)
}

var (
  instanceShorterRepo *shorterRepository
)

// CreateTTLIndex /*
// CreateTTLIndex sets up TTL index on expired_at field
func (r *shorterRepository) CreateTTLIndex(ctx context.Context) error {
  indexModel := mongo.IndexModel{
    Keys:    bson.M{"expired_at": 1},                  //Create index by expired_at
    Options: options.Index().SetExpireAfterSeconds(0), // delete now after URL expire
  }

  //TODO : TTL Index không ngay lập tức xóa dữ liệu, MongoDB chạy một job mỗi 60 giây để kiểm tra và xóa.
  // Đảm bảo expired_at là một MongoDB DateTime, nếu không TTL Index sẽ không hoạt động.

  _, err := r.Collection.Indexes().CreateOne(ctx, indexModel)
  if err != nil {
    return err
  }

  return nil
}

// NewShorterRepository initializes the repository singleton
func NewShorterRepository(dbStorage *mongodb2.DatabaseStorage) IShorterRepository {
  instanceShorterRepo = &shorterRepository{
    Repository: mongodb2.NewRepository[URLs](dbStorage),
  }

  // Create TTL Index
  ctx := context.Background()
  if err := instanceShorterRepo.CreateTTLIndex(ctx); err != nil {
    panic("Failed to create TTL index: " + err.Error())
  }

  return instanceShorterRepo
}

// SaveURL saves the short key and long URL to the database
func (r *shorterRepository) SaveURL(ctx context.Context, shortKey, longURL string, expireIn time.Duration) error {
  now := time.Now()
  url := &URLs{
    ShortKey:  shortKey,
    LongURL:   longURL,
    CreatedAt: now,
    ExpiredAt: now.Add(expireIn),
  }

  _, err := r.Collection.InsertOne(ctx, url)
  if err != nil {
    return err
  }

  return nil
}

// GetURL retrieves the long URL corresponding to the short key, checking expiration
func (r *shorterRepository) GetURL(ctx context.Context, shortKey string) (*URLs, error) {
  var result URLs
  err := r.Collection.FindOne(ctx, bson.M{
    "short_key":  shortKey,
    "expired_at": bson.M{"$gt": time.Now()},
  }).Decode(&result)

  if err != nil {
    return nil, err
  }

  return &result, nil
}
