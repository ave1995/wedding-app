package user

import (
	"context"
	"fmt"
	"wedding-app/domain/model"
	"wedding-app/domain/service"
	"wedding-app/domain/store"

	"github.com/google/uuid"
)

type userService struct {
	store store.UserStore
}

func NewUserService(store store.UserStore) service.UserService {
	return &userService{store: store}
}

func (s *userService) RegisterUser(ctx context.Context, params model.RegisterUserParams) (*model.User, error) {
	return s.store.RegisterUser(ctx, params)
}

// LoginUser implements service.UserService.
func (s *userService) LoginUser(ctx context.Context, email string, password string) (*model.User, error) {
	user, err := s.store.LoginUser(ctx, email, password)
	if err != nil {
		return nil, fmt.Errorf("failed to login user %q: %w", email, err)
	}

	return user, nil
}

// GetUserByID implements service.UserService.
func (s *userService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("failed to parse id %q: %w", id, err)
	}

	return s.store.GetUserByID(ctx, parsed)
}

// CreateGuest implements service.UserService.
func (s *userService) CreateGuest(ctx context.Context, params model.CreateGuestParams) (*model.User, error) {
	if err := uuid.Validate(params.QuizID); err != nil {
		return nil, fmt.Errorf("quiz ID is invalid: %w", err)
	}
	user, err := s.store.CreateGuest(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create guest %q: %w", params.Username, err)
	}

	return user, nil
}
