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
	Update(ctx context.Context, usr User) (User, error)
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
func (u UserServicer) UpdateUser(req UpdateUserRequest, gr server.GenericRequest) UpdateUserResponse {
	if !(gr.Claims.Authorized(RoleAdmin.name) || req.User.ID.String() == gr.Claims.ID) {
		return UpdateUserResponse{Error: fmt.Errorf("Unauthorized action").Error()}
	}
	uu, err := u.storer.Update(gr.Ctx, req.User)
	if err != nil {
		return UpdateUserResponse{Error: err.Error()}
	}
	return UpdateUserResponse{User: uu}

}

// Register implements UserRpcService
func (us UserServicer) Register(s *server.Server) {
	s.Register("UserService", "CreateUser", server.RPCEndpoint{Roles: []string{auth.RoleAdmin}, Handler: us.CreateUserHandler})
	s.Register("UserService", "DeleteUser", server.RPCEndpoint{Roles: []string{auth.RoleAdmin}, Handler: us.DeleteUserHandler})
	s.Register("UserService", "UpdateUser", server.RPCEndpoint{Roles: []string{auth.RoleAdmin, auth.RoleUser}, Handler: us.UpdateUserHandler})
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

// UpdateUserRequest is the request object for UserService.UpdateUser.
type UpdateUserRequest struct {
	User User `json:"user"`
}

// UpdaetUserResponse is the response object for UserService.UpdaetUser.
type UpdateUserResponse struct {
	User  User   `json:"user"`
	Error string `json:"error,omitempty"`
}

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

type QueryUserByIDRequest struct{}
type QueryUserByIDResponse struct{}

type QueryUserByEmailRequest struct{}
type QueryUserByEmailResponse struct{}

type AuthenticateRequest struct{}
type AuthenticateResponse struct{}
