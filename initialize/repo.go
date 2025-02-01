package initialize

import (
  "url-shortener/database/mongodb"
  "url-shortener/repositories/shorter"
)

var (
  repositories *Repositories
)

type Repositories struct {
  // shorter
  shorterRepository shorter.IShorterRepository
}

func NewRepositories(db *mongodb.DatabaseStorage) *Repositories {
  repositories = &Repositories{
    // shorter
    shorterRepository: shorter.NewShorterRepository(db),
  }

  return repositories
}
