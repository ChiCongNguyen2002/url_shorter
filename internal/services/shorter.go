package services

import (
  "context"
  "time"
  "url-shortener/database/redis"
  "url-shortener/repositories/shorter"
  "url-shortener/utils"
)

type ShorterService struct {
  shorterRepo shorter.IShorterRepository
  cache       redis.ICache
}

type IShortenerService interface {
  ShortenURL(c context.Context, longURL string) (string, error)
  RedirectURL(c context.Context, shortKey string) (string, error)
}

func NewShorterService(repo shorter.IShorterRepository, cache redis.ICache) *ShorterService {
  return &ShorterService{
    shorterRepo: repo,
    cache:       cache,
  }
}

func (s ShorterService) ShortenURL(c context.Context, longURL string) (string, error) {
  // Generate unique short key
  var shortener utils.Shortener = utils.Base62Shortener{} // or utils.MD5Shortener{}
  // var shortener utils.Shortener = utils.MD5Shortener{}
  shortKey := shortener.GenerateKey(longURL)

  // Save to repositories
  if err := s.shorterRepo.SaveURL(c, shortKey, longURL); err != nil {
    return "Failed to save URL", err
  }

  // Save to redis
  err := s.cache.Set(c, shortKey, longURL, 24*time.Hour)
  if err != nil {
    return "", err
  }

  return shortKey, nil
}

func (s ShorterService) RedirectURL(c context.Context, shortKey string) (string, error) {
  // Check redis first
  longURL, err := s.cache.Get(c, shortKey)
  if err == nil && longURL != "" {
    return longURL, nil
  }

  // Fetch from repository if not found in redis
  urlObj, err := s.shorterRepo.GetURL(c, shortKey)
  if err != nil {
    return "URL not found", nil
  }

  // Update redis asynchronously (log warning if fails)
  _ = s.cache.Set(c, shortKey, longURL, 24*time.Hour)

  return urlObj.LongURL, nil
}
