package mongodb

type MongoDBConfig struct {
  DatabaseURI  string `env:"DATABASE_URI,required,notEmpty"`
  DatabaseName string `env:"DATABASE_NAME,required,notEmpty"`
}
type MultiConnMongoConfig map[string]map[string]string
