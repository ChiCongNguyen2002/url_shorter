package configs

import (
  "fmt"
  "github.com/caarlos0/env/v7"
  "github.com/joho/godotenv"
  "log"
  "url-shortener/database/mongodb"
  "url-shortener/database/redis"
)

type SystemConfig struct {
  Env           string                `env:"ENV,required,notEmpty"`
  HttpPort      uint64                `env:"HTTP_PORT,required,notEmpty"`
  MongoDBConfig mongodb.MongoDBConfig `envPrefix:"MONGODB_"`
  RedisConfig   redis.RedisConfig     `envPrefix:"REDIS_"`
}

var configSingletonObj *SystemConfig

func LoadConfig() (cf *SystemConfig, err error) {
  if configSingletonObj != nil {
    cf = configSingletonObj
    return
  }

  // Load environment variables from local.env
  err = godotenv.Load("local.env")
  if err != nil {
    fmt.Printf("Warning: Could not load .env file: %v\n\n", err)
  }

  cf = &SystemConfig{}
  if err = env.Parse(cf); err != nil {
    log.Fatalf("failed to unmarshal config: %s", err)
  }

  configSingletonObj = cf
  return
}

func GetInstance() *SystemConfig {
  return configSingletonObj
}
