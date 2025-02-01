package http

import (
  "github.com/labstack/echo/v4"
  "github.com/labstack/echo/v4/middleware"
  "url-shortener/handlers/routers"
)

const healthPath = "/coordinator/api/v1/health"

func (app *Server) InitRouters(e *echo.Echo) error {
  e.Use(middleware.RequestID())
  e.Use(middleware.Recover())
  e.Use(middleware.CORS())
  //e.Use(middlewares.Logging)
  //e.Use(middlewares.AddExtraDataForRequestContext)
  //e.Use(middlewares.Region)
  //e.GET(healthPath, func(c echo.Context) error {
  //  if healthCheck {
  //    return c.JSON(http.StatusOK, resp.BuildSuccessResp(resp.LangEN, nil))
  //  }
  //
  //  return c.JSON(http.StatusInternalServerError, resp.BuildErrorResp(resp.ErrSystem, "", resp.LangEN))
  //})

  // SHORTER
  controller := routers.NewControllers(e, app.handlers)
  controller.SetupRoutes()
  return nil
}
