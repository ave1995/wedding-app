package service

import (
	"context"
	"wedding-app/domain/model"
)

type UserService interface {
	RegisterUser(ctx context.Context, params model.RegisterUserParams) (*model.User, error)
	LoginUser(ctx context.Context, email, password string) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	CreateGuest(ctx context.Context, params model.CreateGuestParams) (*model.User, error)
}
