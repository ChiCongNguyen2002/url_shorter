package main

import (
  "context"
  "fmt"
  "github.com/labstack/echo/v4"
  "os"
  "os/signal"
  "syscall"
  "time"
  "url-shortener/configs"
  "url-shortener/database/redis"
  http "url-shortener/handlers"
  "url-shortener/initialize"
  "url-shortener/logger"
)

func main() {
  // Load configuration
  conf, err := configs.LoadConfig()
  if err != nil {
    //log.Fatal().Msgf("Load config failed! %s", err)
    fmt.Printf("Load config failed! %s\n", err)
  }

  logger.InitLog(os.Getenv("SERVICE_ID"))
  log := logger.GetLogger()
  log.Info().Any("service", os.Getenv("SERVICE_ID")).Msg("Start services")
  // Set up Echo
  e := echo.New()

  http.SetHealthCheck(true)
  ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt, syscall.SIGTERM)
  defer cancel()

  // Initialize database connections
  databaseStorage := initialize.NewDatabaseConnection(ctx)

  // Initialize repositories
  repo := initialize.NewRepositories(databaseStorage.Conn)

  // Initialize Cache
  cacheService := redis.NewRedisCache(&conf.RedisConfig)

  // Initialize services
  services := initialize.NewServices(repo, cacheService)

  // Initialize handlers
  handlers := initialize.NewHandlers(services)

  // Start HTTP server
  srv := http.NewHttpServe(handlers)
  srv.Start(e)

  <-ctx.Done()
  http.SetHealthCheck(false)
  shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
  defer shutdownCancel()
  if err := e.Shutdown(shutdownCtx); err != nil {

  }
}
