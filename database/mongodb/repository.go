package mongodb

import (
  "context"
  "fmt"
  "log"
  "time"

  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
)

// ModelInterface defines methods for retrieving collection names and indexes
type ModelInterface interface {
  CollectionName() string
  IndexModels() []mongo.IndexModel
}

// Repository provides generic database operations
type Repository[T ModelInterface] struct {
  *mongo.Collection
}

// NewRepository initializes a repository
func NewRepository[T ModelInterface](dbStorage *DatabaseStorage, opts ...*options.CollectionOptions) *Repository[T] {
  var t T
  collectionName := t.CollectionName()
  indexModels := t.IndexModels()

  if dbStorage.db == nil {
    log.Fatalf("database instance is nil")
  }

  collection, err := newRepository(dbStorage.db, collectionName, indexModels, opts...)
  if err != nil {
    log.Fatalf("new repository error: %v", err)
  }

  return &Repository[T]{
    Collection: collection,
  }
}

func newRepository(db *mongo.Database, collectionName string, indexModels []mongo.IndexModel, opts ...*options.CollectionOptions) (*mongo.Collection, error) {
  collection := db.Collection(collectionName, opts...)

  if len(indexModels) > 0 {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    if _, err := collection.Indexes().CreateMany(ctx, indexModels); err != nil {
      return nil, fmt.Errorf("create index collectionName=%v error: %v", collectionName, err)
    }
  }

  return collection, nil
}
