package initialize

import (
  "url-shortener/handlers/api"
)

type Handlers struct {
  ShorterHandler *api.ShorterHandler
}

func NewHandlers(services *Services) *Handlers {
  shorterHandler := api.NewShorterHandler(
    services.shorterService,
  )

  return &Handlers{
    ShorterHandler: shorterHandler,
  }
}
