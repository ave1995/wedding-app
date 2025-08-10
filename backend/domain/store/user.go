package store

import (
	"context"
	"wedding-app/domain/model"

	"github.com/google/uuid"
)

type UserStore interface {
	RegisterUser(ctx context.Context, params model.RegisterUserParams) (*model.User, error)
	LoginUser(ctx context.Context, email, password string) (*model.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	CreateGuest(ctx context.Context, params model.CreateGuestParams) (*model.User, error)
}
