package services

import (
	"context"
	"wedding-app/models/domain"
	"wedding-app/repositories"
)

type UserService interface {
	RegisterUser(ctx context.Context, username, email, password string) (*domain.User, error)
}

type userService struct {
	repository repositories.UserRepository
}

func NewUserService(repository repositories.UserRepository) UserService {
	return &userService{repository: repository}
}

func (s *userService) RegisterUser(ctx context.Context, username, email, password string) (*domain.User, error) {
	return s.repository.RegisterUser(ctx, username, email, password)
}
