package initialize

import (
  "url-shortener/database/redis"
  "url-shortener/internal/services"
)

type Services struct {
  shorterService services.IShortenerService
}

func NewServices(repo *Repositories, cache redis.ICache) *Services {
  shorterService := services.NewShorterService(
    repo.shorterRepository,
    cache,
  )
  service := &Services{
    shorterService: shorterService,
  }
  return service
}
