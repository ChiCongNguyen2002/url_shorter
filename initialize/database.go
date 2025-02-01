package initialize

import (
  "context"
  "fmt"
  "url-shortener/configs"
  "url-shortener/database/mongodb"
)

var (
  databaseConnection *DatabaseConnection
)

type DatabaseConnection struct {
  Conn *mongodb.DatabaseStorage
}

func NewDatabaseConnection(ctx context.Context) *DatabaseConnection {
  test := &configs.GetInstance().MongoDBConfig
  fmt.Println("test", test)
  conn, err := mongodb.ConnectMongoDB(ctx, &configs.GetInstance().MongoDBConfig)
  if err != nil {
    fmt.Println("Dont Connect Database")
  }
  databaseConnection = &DatabaseConnection{
    Conn: conn,
  }

  return databaseConnection
}
