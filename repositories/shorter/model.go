package shorter

import (
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  "time"
)

type URLs struct {
  ShortKey  string    `bson:"short_key"`
  LongURL   string    `bson:"long_url"`
  CreatedAt time.Time `bson:"created_at"`
  ExpiredAt time.Time `bson:"expired_at"`
}

func (l URLs) IndexModels() []mongo.IndexModel {
  return []mongo.IndexModel{
    {
      Keys: bson.D{
        {Key: "short_key", Value: 1}, // Index on 'short_key' field
      },
      Options: options.Index().SetUnique(true),
    },
  }
}

func (URLs) CollectionName() string {
  return "urls"
}
