package user

import (
	"context"
	"fmt"
	"wedding-app/domain/model"
	"wedding-app/domain/service"
	"wedding-app/domain/store"
)

type userService struct {
	store      store.UserStore
	jwtService service.JWTService
}

func NewUserService(store store.UserStore, jwtService service.JWTService) service.UserService {
	return &userService{store: store, jwtService: jwtService}
}

func (s *userService) RegisterUser(ctx context.Context, username, email, password string) (*model.User, error) {
	return s.store.RegisterUser(ctx, username, email, password)
}

// LoginUser implements service.UserService.
func (s *userService) LoginUser(ctx context.Context, email string, password string) (*model.AccessToken, error) {
	user, err := s.store.LoginUser(ctx, email, password)
	if err != nil {
		return nil, fmt.Errorf("failed to login user %q: %w", email, err)
	}

	return s.jwtService.Generate(user)
}
