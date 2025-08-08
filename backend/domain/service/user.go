package service

import (
	"context"
	"wedding-app/domain/model"
)

type UserService interface {
	RegisterUser(ctx context.Context, username, email, password string) (*model.User, error)
	LoginUser(ctx context.Context, email, password string) (*model.AccessToken, error)
}
