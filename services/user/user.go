package user

import (
	"context"
	"fmt"
	"net/mail"

	"github.com/gitamped/seed/auth"
	"github.com/gitamped/seed/server"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
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
	Create(ctx context.Context, usr User) (User, error)
	Delete(ctx context.Context, email mail.Address) (User, error)
	QueryByID(ctx context.Context, id string) (User, error)
	QueryByEmail(ctx context.Context, email string) (User, error)
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
func (u UserServicer) QueryUserByEmail(req QueryUserByEmailRequest, gr server.GenericRequest) QueryUserByEmailResponse {
	usr, err := u.storer.QueryByEmail(gr.Ctx, req.Email)
	if err != nil {
		return QueryUserByEmailResponse{Error: err.Error()}
	}
	return QueryUserByEmailResponse{User: usr}
}

// QueryUserByID implements UserRpcService
func (u UserServicer) QueryUserByID(req QueryUserByIDRequest, gr server.GenericRequest) QueryUserByIDResponse {
	usr, err := u.storer.QueryByID(gr.Ctx, req.ID)
	if err != nil {
		return QueryUserByIDResponse{Error: err.Error()}
	}
	return QueryUserByIDResponse{User: usr}
}

// QueryUser implements UserRpcService
func (UserServicer) QueryUser(QueryUserRequest, server.GenericRequest) QueryUserResponse {
	panic("unimplemented")
}

// DeleteUser implements UserRpcService
func (u UserServicer) DeleteUser(req DeleteUserRequest, gr server.GenericRequest) DeleteUserResponse {
	du, err := u.storer.Delete(gr.Ctx, req.User.Email)
	if err != nil {
		return DeleteUserResponse{Error: err.Error()}
	}
	return DeleteUserResponse{User: du}
}

// CreateUser implements UserRpcService
func (u UserServicer) CreateUser(req CreateUserRequest, gr server.GenericRequest) CreateUserResponse {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return CreateUserResponse{Error: fmt.Errorf("generatefrompassword: %w", err).Error()}
	}
	usr := User{
		ID:           uuid.New(),
		Name:         req.NewUser.Name,
		Email:        req.NewUser.Email,
		PasswordHash: hash,
		Roles:        req.NewUser.Roles,
		Department:   req.NewUser.Department,
		Enabled:      true,
		DateCreated:  gr.Values.Now,
		DateUpdated:  gr.Values.Now,
	}
	result, err := u.storer.Create(gr.Ctx, usr)
	if err != nil {
		return CreateUserResponse{Error: err.Error()}
	}
	return CreateUserResponse{User: result}
}

// UpdateUser implements UserRpcService
func (UserServicer) UpdateUser(UpdateUserRequest, server.GenericRequest) UpdateUserResponse {
	panic("unimplemented")
}

// Register implements UserRpcService
func (us UserServicer) Register(s *server.Server) {
	s.Register("UserService", "CreateUser", server.RPCEndpoint{Roles: []string{auth.RoleAdmin}, Handler: us.CreateUserHandler})
	s.Register("UserService", "DeleteUser", server.RPCEndpoint{Roles: []string{auth.RoleAdmin}, Handler: us.DeleteUserHandler})
	s.Register("UserService", "QueryUserByID", server.RPCEndpoint{Roles: []string{auth.RoleAdmin}, Handler: us.QueryUserByIDHandler})
	s.Register("UserService", "QueryUserByEmail", server.RPCEndpoint{Roles: []string{auth.RoleAdmin}, Handler: us.QueryUserByEmailHandler})
}

// Create new UserServicer
func NewUserServicer(log *zap.SugaredLogger, storer Storer) UserRpcService {
	return UserServicer{
		log:    log,
		storer: storer,
	}
}

// CreateUserRequest is the request object for UserService.CreateUser.
type CreateUserRequest struct {
	NewUser NewUser `json:"newUser"`
}

// CreateUserResponse is the response object containing a UserService.CreateUser.
type CreateUserResponse struct {
	User  User   `json:"user"`
	Error string `json:"error,omitempty"`
}

type UpdateUserRequest struct{}
type UpdateUserResponse struct{}

// DeleteUserRequest is the request object for UserService.DeleteUser.
type DeleteUserRequest struct {
	User User `json:"user"`
}

// DeleteUserResponse is the response object for UserService.DeleteUser.
type DeleteUserResponse struct {
	User  User   `json:"user"`
	Error string `json:"error,omitempty"`
}

type QueryUserRequest struct{}
type QueryUserResponse struct{}

// QueryUserByIDRequest is the request object for UserService.QueryUserByID.
type QueryUserByIDRequest struct {
	ID string `json:"id"`
}

// QueryUserByIDResponse is the response object for UserService.QueryUserByID.
type QueryUserByIDResponse struct {
	User  User   `json:"user"`
	Error string `json:"error,omitempty"`
}

// QueryUserByEmailRequest is the request object for UserService.QueryUserByEmail.
type QueryUserByEmailRequest struct {
	Email string `json:"email"`
}

// QueryUserByEmailResponse is the response object for UserService.QueryUserByEmail.
type QueryUserByEmailResponse struct {
	User  User   `json:"user"`
	Error string `json:"error,omitempty"`
}

type AuthenticateRequest struct{}
type AuthenticateResponse struct{}
