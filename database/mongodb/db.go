package mongodb

import (
  "context"
  "encoding/base64"
  "encoding/json"
  "fmt"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  "go.mongodb.org/mongo-driver/mongo/readpref"
  "go.mongodb.org/mongo-driver/mongo/writeconcern"
  "log"
  "strings"
  "time"
  "url-shortener/utils"
)

type DatabaseStorage struct {
  db            *mongo.Database
  client        *mongo.Client
  mappingDB     map[string]*mongo.Database
  mappingClient map[string]*mongo.Client
}

type SessionMultiConn struct {
  clients map[string]*mongo.Client
}

var dbStorage *DatabaseStorage

type RepoTxMultiConnInterface interface {
  CollectionName() string
}

func ConnectMongoDB(ctx context.Context, config *MongoDBConfig, multiConnCfg ...string) (*DatabaseStorage, error) {
  if dbStorage != nil {
    return dbStorage, nil
  }

  if config != nil {
    client, db, err := connect(ctx, config)
    if err != nil {
      return nil, err
    }

    dbStorage = &DatabaseStorage{
      db:     db,
      client: client,
    }
    return dbStorage, nil
  }

  if len(multiConnCfg) == 0 {
    return nil, fmt.Errorf("multi conn config not found for mongodb")
  }

  cfgByte, err := base64.StdEncoding.DecodeString(multiConnCfg[0])
  if err != nil {
    return nil, fmt.Errorf("base64 decode multi conn config failed: %v", err)
  }

  var multiCfg MultiConnMongoConfig
  err = json.Unmarshal(cfgByte, &multiCfg)
  if err != nil {
    return nil, fmt.Errorf("unmarshal multi conn config failed: %v", err)
  }

  mappingDB := make(map[string]*mongo.Database)
  mappingClient := make(map[string]*mongo.Client)

  for region, value := range multiCfg {
    for dbName, uri := range value {
      connName := fmt.Sprintf("%s::%s", region, dbName)
      cfg := &MongoDBConfig{
        DatabaseURI:  uri,
        DatabaseName: dbName,
      }

      client, db, err := connect(ctx, cfg)
      if err != nil {
        return nil, err
      }

      mappingDB[connName] = db
      mappingClient[connName] = client
    }
  }

  dbStorage = &DatabaseStorage{
    mappingDB:     mappingDB,
    mappingClient: mappingClient,
  }

  return dbStorage, nil
}

func connect(ctx context.Context, config *MongoDBConfig) (*mongo.Client, *mongo.Database, error) {
  // Set a timeout for the connection
  ctxNew, cancel := context.WithTimeout(ctx, 30*time.Second)
  defer cancel()

  // Set client options
  clientOpts := options.Client().ApplyURI(config.DatabaseURI)
  clientOpts.SetMaxPoolSize(100)

  // Connect to MongoDB
  client, err := mongo.Connect(ctxNew, clientOpts)
  if err != nil {
    log.Printf("Failed to create MongoDB client: %v", err)
    return nil, nil, err
  }

  // Verify connection with a ping
  if err := client.Ping(ctxNew, readpref.Primary()); err != nil {
    log.Printf("Failed to ping MongoDB: %v", err)
    return nil, nil, err
  }

  log.Println("Successfully connected to MongoDB")

  db := client.Database(config.DatabaseName)
  return client, db, nil
}

func GetDatabaseStorage() *DatabaseStorage {
  return dbStorage
}

func (dbStorage *DatabaseStorage) StartSession(opts ...*options.SessionOptions) (mongo.Session, error) {
  return dbStorage.client.StartSession(opts...)
}

func (dbStorage *DatabaseStorage) GetClient() *mongo.Client {
  return dbStorage.client
}

func (dbStorage *DatabaseStorage) StartSessionMultiConn(ctx context.Context, opts ...*options.SessionOptions) (mongo.Session, error) {
  if dbStorage.mappingClient == nil {
    return nil, fmt.Errorf("mongo multi conn: mappingClient nil pointer")
  }

  connName, ok := ctx.Value(utils.KeyMongoMultiConnName).(string)
  if !ok {
    return nil, fmt.Errorf("mongo_multi_conn_name not found in context")
  }

  client, ok := dbStorage.mappingClient[connName]
  if !ok {
    return nil, fmt.Errorf("mongo multi conn: client not found: conn_name=%s", connName)
  }

  return client.StartSession(opts...)
}

func (dbStorage *DatabaseStorage) GetClientMultiConn(ctx context.Context) (*mongo.Client, error) {
  if dbStorage.mappingClient == nil {
    return nil, fmt.Errorf("mongo multi conn: mappingClient nil pointer")
  }

  connName, ok := ctx.Value(utils.KeyMongoMultiConnName).(string)
  if !ok {
    return nil, fmt.Errorf("mongo_multi_conn_name not found in context")
  }

  client, ok := dbStorage.mappingClient[connName]
  if !ok {
    return nil, fmt.Errorf("mongo multi conn: client not found: conn_name=%s", connName)
  }

  return client, nil
}

func (dbStorage *DatabaseStorage) ExecTransaction(ctx context.Context, callback func(sessCtx mongo.SessionContext) (interface{}, error)) error {
  if dbStorage.client == nil {
    return fmt.Errorf("client nil pointer")
  }

  // start session
  session, err := dbStorage.client.StartSession()
  if err != nil {
    return err
  }

  // end session
  defer session.EndSession(ctx)

  wc := writeconcern.Majority()
  opts := options.Transaction().SetWriteConcern(wc)

  // execute transaction
  _, err = session.WithTransaction(ctx, callback, opts)
  return err
}

func (dbStorage *DatabaseStorage) InitSessionMultiConn(dbNames ...string) (*SessionMultiConn, error) {
  if len(dbNames) == 0 {
    return nil, fmt.Errorf("InitTxMultiConn: dbNames empty")
  }

  if dbStorage.mappingClient == nil {
    return nil, fmt.Errorf("InitTxMultiConn: mappingClient nil pointer")
  }

  mapDBNames := make(map[string]bool)
  for _, dbName := range dbNames {
    mapDBNames[dbName] = true
  }

  clients := make(map[string]*mongo.Client)
  for connName, client := range dbStorage.mappingClient {
    split := strings.Split(connName, "::")
    if len(split) != 2 {
      return nil, fmt.Errorf("InitTxMultiConn: split connName failed: connName=%s", connName)
    }

    if _, ok := mapDBNames[split[1]]; ok {
      clients[split[0]] = client
    }
  }

  return &SessionMultiConn{clients: clients}, nil
}
