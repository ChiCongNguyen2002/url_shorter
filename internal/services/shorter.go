package services

import (
  "context"
  "fmt"
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
  ShortenURL(c context.Context, longURL string, expiredAt time.Duration) (string, error)
  RedirectURL(c context.Context, shortKey string) (string, error)
}

func NewShorterService(repo shorter.IShorterRepository, cache redis.ICache) *ShorterService {
  return &ShorterService{
    shorterRepo: repo,
    cache:       cache,
  }
}

func (s ShorterService) ShortenURL(c context.Context, longURL string, expiredAt time.Duration) (string, error) {
  // Generate unique short key
  var shortener utils.Shortener = utils.Base62Shortener{} // or utils.MD5Shortener{}
  // var shortener utils.Shortener = utils.MD5Shortener{}
  shortKey := shortener.GenerateKey(longURL)

  // Check URL exist
  existingURL, err := s.cache.Get(c, shortKey)
  if err == nil && existingURL != "" {
    fmt.Printf("🔄 URL already shortened: %s -> %s", longURL, shortKey)
    return shortKey, nil
  }

  // Save to repositories
  if err := s.shorterRepo.SaveURL(c, shortKey, longURL, expiredAt); err != nil {
    return "Failed to save URL", err
  }

  // Save to redis
  if err := s.cache.Set(c, shortKey, longURL, expiredAt); err != nil {
    fmt.Printf("Warning: Failed to cache URL %s: %v", shortKey, err)
  }

  return shortKey, nil
}

func (s ShorterService) RedirectURL(c context.Context, shortKey string) (string, error) {
  // Check redis first
  longURL, err := s.cache.Get(c, shortKey)
  if err == nil && longURL != "" {
    return longURL, err
  }

  // Fetch from repository if not found in redis
  urlObj, err := s.shorterRepo.GetURL(c, shortKey)
  if err != nil {
    return "URL not found", err
  }

  // Update redis asynchronously (log warning if fails)
  go func() {
    if cacheErr := s.cache.Set(c, shortKey, urlObj.LongURL, 24*time.Hour); cacheErr != nil {
      fmt.Printf("⚠️ Warning: Failed to cache URL %s: %v", shortKey, cacheErr)
    }
  }()

  return urlObj.LongURL, nil
}
