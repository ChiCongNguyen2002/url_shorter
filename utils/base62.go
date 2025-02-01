package utils

import (
  "time"
)

type Base62Shortener struct{}

func (b Base62Shortener) GenerateKey(input string) string {
  return EncodeBase62(time.Now().UnixNano())
}

func EncodeBase62(num int64) string {
  const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
  result := ""
  for num > 0 {
    remainder := num % 62
    result = string(base62Chars[remainder]) + result
    num /= 62
  }
  return result
}
