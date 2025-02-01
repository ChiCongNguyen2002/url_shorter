package shorter

import (
  "context"
  "go.mongodb.org/mongo-driver/bson"
  "time"
  mongodb2 "url-shortener/database/mongodb"
)

type shorterRepository struct {
  *mongodb2.Repository[URLs]
}

type IShorterRepository interface {
  SaveURL(ctx context.Context, shortKey, longURL string) error
  GetURL(ctx context.Context, shortKey string) (*URLs, error)
}

var (
  instanceShorterRepo *shorterRepository
)

// NewShorterRepository initializes the repository singleton
func NewShorterRepository(dbStorage *mongodb2.DatabaseStorage) IShorterRepository {
  instanceShorterRepo = &shorterRepository{
    Repository: mongodb2.NewRepository[URLs](dbStorage),
  }
  return instanceShorterRepo
}

// SaveURL saves the short key and long URL to the database
func (r *shorterRepository) SaveURL(ctx context.Context, shortKey, longURL string) error {
  
  url := &URLs{
    ShortKey:  shortKey,
    LongURL:   longURL,
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
  }

  _, err := r.Collection.InsertOne(ctx, url)
  if err != nil {
    return err
  }

  return nil
}

// GetURL retrieves the long URL corresponding to the short key
func (r *shorterRepository) GetURL(ctx context.Context, shortKey string) (*URLs, error) {
  var result URLs
  err := r.Collection.FindOne(ctx, bson.M{"short_key": shortKey}).Decode(&result)
  if err != nil {
    return nil, err
  }
  return &result, nil
}
