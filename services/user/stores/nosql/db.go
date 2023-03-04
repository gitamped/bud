package db

import (
	"github.com/arangodb/go-driver"
	"go.uber.org/zap"
)

type Store struct {
	db  driver.Database
	log *zap.SugaredLogger
}

// NewStore constructs the api for data access.
func NewStore(log *zap.SugaredLogger, db driver.Database) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}
