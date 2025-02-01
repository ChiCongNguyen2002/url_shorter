package http

import (
  "errors"
  "fmt"
  "net/http"
  "sync"
  "url-shortener/configs"
  "url-shortener/initialize"

  "github.com/labstack/echo/v4"
)

var (
  healthCheck bool
  mu          sync.RWMutex
)

func SetHealthCheck(status bool) {
  mu.Lock()
  defer mu.Unlock()
  healthCheck = status
}

type ServInterface interface {
  Start(e *echo.Echo)
}

type Server struct {
  handlers *initialize.Handlers
}

func NewHttpServe(handlers *initialize.Handlers,
) *Server {
  return &Server{
    handlers: handlers,
  }
}

func (app *Server) Start(e *echo.Echo) {
  err := app.InitRouters(e)
  if err != nil {
    //log.Fatal().Msgf("InitRouters fail! %s", err)
  }

  httpPort := configs.GetInstance().HttpPort
  go func() {
    err := e.Start(fmt.Sprintf(":%d", httpPort))
    if err != nil && !errors.Is(err, http.ErrServerClosed) {
      //log.Fatal().Msgf("can't start echo")
    }
  }()
  // log.Info().Msg("all services already")
}
