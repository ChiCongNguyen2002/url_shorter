package api

import (
  "github.com/labstack/echo/v4"
  "net/http"
  "url-shortener/internal/services"
  "url-shortener/models"
  "url-shortener/utils"
)

type ShorterHandler struct {
  shorterHandler services.IShortenerService
}

func NewShorterHandler(shorterHandler services.IShortenerService) (handler *ShorterHandler) {
  return &ShorterHandler{
    shorterHandler: shorterHandler,
  }
}

func (s ShorterHandler) ShortenURL(c echo.Context) error {
  var reqLongURL models.RequestBody
  if err := c.Bind(&reqLongURL); err != nil {
    return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
  }

  shortKey, err := s.shorterHandler.ShortenURL(c.Request().Context(), reqLongURL.LongURL, utils.TimeExpireURL)
  if err != nil {
    return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save URL"})
  }

  return c.JSON(http.StatusOK, map[string]string{"short_url": "http://localhost:3001/api/v1/" + shortKey})
}

func (s ShorterHandler) RedirectURL(c echo.Context) error {
  shortKey := c.Param("short_key")
  if shortKey == "" {
    return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
  }
  longURL, err := s.shorterHandler.RedirectURL(c.Request().Context(), shortKey)
  if err != nil {
    return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save URL"})
  }
  return c.Redirect(http.StatusMovedPermanently, longURL)
}
