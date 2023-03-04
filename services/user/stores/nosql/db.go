package nosql

import (
	"context"

	"github.com/arangodb/go-driver"
	"github.com/gitamped/bud/services/user"
	"go.uber.org/zap"
)

type Store struct {
	db  driver.Database
	col driver.Collection
	log *zap.SugaredLogger
}

// NewStore constructs the api for data access.
func NewStore(log *zap.SugaredLogger, db driver.Database) *Store {
	col, err := db.Collection(context.Background(), "users")
	if err != nil {
		log.Panicf("error accessing collection: %s", err)
	}
	return &Store{
		log: log,
		db:  db,
		col: col,
	}
}

// Create inserts a new user into the database.
func (s *Store) Create(ctx context.Context, usr user.User) error {
	var result dbUser
	driver.WithReturnNew(ctx, result)
	_, err := s.col.CreateDocument(ctx, usr)
	return err
}
