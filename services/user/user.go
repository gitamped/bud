package user

import (
	"github.com/gitamped/seed/auth"
	"github.com/gitamped/seed/server"
)

// UserService is an API for creating users for an app.
type UserService interface {
	// CreateUser create a user
	CreateUser(CreateUserRequest, server.GenericRequest) CreateUserResponse
	// UpdateUser updates a user
	UpdateUser(UpdateUserRequest, server.GenericRequest) UpdateUserResponse
	// DeleteUser deletes a user
	DeleteUser(DeleteUserRequest, server.Server) DeleteUserResponse
}

// Required to register endpoints with the Server
type UserRpcService interface {
	UserService
	// Registers RPCService with Server
	Register(s *server.Server)
}

// Implements interface
type UserServicer struct{}

// DeleteUser implements UserRpcService
func (UserServicer) DeleteUser(DeleteUserRequest, server.Server) DeleteUserResponse {
	panic("unimplemented")
}

// CreateUser implements UserRpcService
func (UserServicer) CreateUser(req CreateUserRequest, gr server.GenericRequest) CreateUserResponse {
	// TODO call db layer
	u := CreateUserResponse{}
	u.Name = "John Doe"
	return u
}

// UpdateUser implements UserRpcService
func (UserServicer) UpdateUser(UpdateUserRequest, server.GenericRequest) UpdateUserResponse {
	panic("unimplemented")
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

type UpdateUserRequest struct{}
type UpdateUserResponse struct{}

type DeleteUserRequest struct{}
type DeleteUserResponse struct{}
