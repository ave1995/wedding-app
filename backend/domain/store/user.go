package store

import (
	"context"
	"wedding-app/domain/model"
)

type UserStore interface {
	RegisterUser(ctx context.Context, username, email, password string) (*model.User, error)
	LoginUser(ctx context.Context, email, password string) (*model.User, error)
}
