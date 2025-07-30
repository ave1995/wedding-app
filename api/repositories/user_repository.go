package repositories

import (
	"context"
	"wedding-app/models/domain"
)

type UserRepository interface {
	RegisterUser(ctx context.Context, username, email, password string) (*domain.User, error)
}
