package routers

import (
  "github.com/labstack/echo/v4"
  "url-shortener/initialize"
)

type Controllers struct {
  e         *echo.Echo
  clientSys *echo.Group
  handlers  *initialize.Handlers
}

func NewControllers(e *echo.Echo, handlers *initialize.Handlers) *Controllers {
  return &Controllers{
    e:         e,
    clientSys: e.Group(prefixSystemPath),
    handlers:  handlers,
  }
}

func (app *Controllers) SetupRoutes() {
  // Set up router for shorter
  app.SetupShorterRouters()
}

func (app *Controllers) SetupShorterRouters() {
  profile := app.clientSys.Group(prefixSystemShorterVersionPath)
  profile.GET(prefixRedirectShortKeyPath, app.handlers.ShorterHandler.RedirectURL)
  profile.POST(prefixLongURLPath, app.handlers.ShorterHandler.ShortenURL)
}
