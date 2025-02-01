package utils

type Shortener interface {
  GenerateKey(input string) string
}
