package nosql

import (
	"context"
	"net/mail"

	"github.com/arangodb/go-driver"
	"github.com/gitamped/bud/services/user"
	"go.uber.org/zap"
)

const collectionName = "users"

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
func (s *Store) Create(ctx context.Context, usr user.User) (user.User, error) {
	var result dbUser
	ctx = driver.WithReturnNew(ctx, &result)
	_, err := s.col.CreateDocument(ctx, toDBUser(usr))
	return toCoreUser(result), err
}

// Delete deletes a user from the database
func (s *Store) Delete(ctx context.Context, email mail.Address) (user.User, error) {
	var result dbUser
	ctx = driver.WithReturnOld(ctx, &result)
	_, err := s.col.RemoveDocument(ctx, email.Address)
	return toCoreUser(result), err
}

// Update updates a user in the database.
func (s *Store) Update(ctx context.Context, usr user.User) (user.User, error) {
	var result dbUser
	ctx = driver.WithReturnNew(ctx, &result)
	_, err := s.col.UpdateDocument(ctx, usr.Email.Address, toDBUser(usr))
	return toCoreUser(result), err
}
