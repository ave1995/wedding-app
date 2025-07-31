package user

import (
	"context"
	"wedding-app/domain/model"
	"wedding-app/domain/service"
	"wedding-app/domain/store"
)

type userService struct {
	store store.UserStore
}

func NewUserService(store store.UserStore) service.UserService {
	return &userService{store: store}
}

func (s *userService) RegisterUser(ctx context.Context, username, email, password string) (*model.User, error) {
	return s.store.RegisterUser(ctx, username, email, password)
}
