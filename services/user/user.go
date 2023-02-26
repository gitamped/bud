package user

import (
	"encoding/json"
	"fmt"

	"github.com/gitamped/seed/auth"
	"github.com/gitamped/seed/server"
	"github.com/gitamped/seed/validate"
)

// UserService is an API for creating users for an app.
type UserService interface {
	// CreateUser create a user
	CreateUser(CreateUserRequest, server.GenericRequest) CreateUserResponse
}

// Required to register endpoints with the Server
type UserRpcService interface {
	UserService
	// Registers RPCService with Server
	Register(s *server.Server)
}

// Implements interface
type UserServicer struct{}

// CreateUserHandler validates input data prior to calling CreateUser
func (us UserServicer) CreateUserHandler(r server.GenericRequest, b []byte) (any, error) {
	var ur CreateUserRequest
	if err := json.Unmarshal(b, &ur); err != nil {
		return nil, fmt.Errorf("Unmarshalling data: %w", err)
	}

	if err := validate.Check(ur); err != nil {
		return nil, fmt.Errorf("validating data: %w", err)
	}

	return us.CreateUser(ur, r), nil
}

// CreateUser implements UserRpcService
func (UserServicer) CreateUser(req CreateUserRequest, gr server.GenericRequest) CreateUserResponse {
	// TODO call db layer
	u := CreateUserResponse{}
	u.Name = "Gopher"
	return u
}

// Register implements UserRpcService
func (us UserServicer) Register(s *server.Server) {
	s.Register("UserService", "CreateUser", server.RPCEndpoint{Roles: []string{auth.RoleAdmin}, Handler: us.CreateUserHandler})
}

// Create new UserServicer
func NewUserServicer() UserRpcService {
	return UserServicer{}
}

// CreateUserRequest is the request object for UserService.Greet.
type CreateUserRequest struct {
	NewUser
}

// CreateUserResponse is the response object containing a
// person's greeting.
type CreateUserResponse struct {
	User
}
