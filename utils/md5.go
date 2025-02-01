package utils

import (
  "crypto/md5"
  "encoding/hex"
)

// MD5Shortener sử dụng MD5 hashing
type MD5Shortener struct{}

// GenerateKey cho MD5
func (m MD5Shortener) GenerateKey(input string) string {
  hash := md5.Sum([]byte(input))
  return hex.EncodeToString(hash[:])[:8]
}
