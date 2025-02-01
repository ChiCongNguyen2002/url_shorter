package models

import "time"

type RequestBody struct {
  LongURL string `json:"long_url"`
}

type URLMapping struct {
  ShortKey  string    `bson:"short_key" json:"short_key"`
  LongURL   string    `bson:"long_url" json:"long_url"`
  CreatedAt time.Time `bson:"created_at" json:"created_at"`
  UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
