package user

import (
	"context"

	"github.com/gitamped/seed/auth"
	"github.com/gitamped/seed/server"
	"go.uber.org/zap"
)

// UserService is an API for creating users for an app.
type UserService interface {
	// CreateUser create a user
	CreateUser(CreateUserRequest, server.GenericRequest) CreateUserResponse
	// UpdateUser updates a user
	UpdateUser(UpdateUserRequest, server.GenericRequest) UpdateUserResponse
	// DeleteUser deletes a user
	DeleteUser(DeleteUserRequest, server.GenericRequest) DeleteUserResponse
	// QueryUser retrieves a list of existing users
	QueryUser(QueryUserRequest, server.GenericRequest) QueryUserResponse
	// QueryByID gets the specified user by id
	QueryUserByID(QueryUserByIDRequest, server.GenericRequest) QueryUserByIDResponse
	// QueryByEmail gets the specified user by email
	QueryUserByEmail(QueryUserByEmailRequest, server.GenericRequest) QueryUserByEmailResponse
	// Authenticate finds a user by their email and verifies their password. On
	// success it returns a Claims User representing this user. The claims can be
	// used to generate a token for future authentication.
	Authenticate(AuthenticateRequest, server.GenericRequest) AuthenticateResponse
}

// Storer interface declares the behavior this package needs to perists and
// retrieve data.
type Storer interface {
	Create(ctx context.Context, usr User) error
}

// Required to register endpoints with the Server
type UserRpcService interface {
	UserService
	// Registers RPCService with Server
	Register(s *server.Server)
}

// Implements interface
type UserServicer struct {
	log    *zap.SugaredLogger
	storer Storer
}

// Authenticate implements UserRpcService
func (UserServicer) Authenticate(AuthenticateRequest, server.GenericRequest) AuthenticateResponse {
	panic("unimplemented")
}

// QueryUserByEmail implements UserRpcService
func (UserServicer) QueryUserByEmail(QueryUserByEmailRequest, server.GenericRequest) QueryUserByEmailResponse {
	panic("unimplemented")
}

// QueryUserByID implements UserRpcService
func (UserServicer) QueryUserByID(QueryUserByIDRequest, server.GenericRequest) QueryUserByIDResponse {
	panic("unimplemented")
}

// QueryUser implements UserRpcService
func (UserServicer) QueryUser(QueryUserRequest, server.GenericRequest) QueryUserResponse {
	panic("unimplemented")
}

// DeleteUser implements UserRpcService
func (UserServicer) DeleteUser(DeleteUserRequest, server.GenericRequest) DeleteUserResponse {
	panic("unimplemented")
}

// CreateUser implements UserRpcService
func (u UserServicer) CreateUser(req CreateUserRequest, gr server.GenericRequest) CreateUserResponse {
	panic("uninplemented")
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
func NewUserServicer(log *zap.SugaredLogger, storer Storer) UserRpcService {
	return UserServicer{
		log:    log,
		storer: storer,
	}
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

type QueryUserRequest struct{}
type QueryUserResponse struct{}

type QueryUserByIDRequest struct{}
type QueryUserByIDResponse struct{}

type QueryUserByEmailRequest struct{}
type QueryUserByEmailResponse struct{}

type AuthenticateRequest struct{}
type AuthenticateResponse struct{}
