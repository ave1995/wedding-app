package service

import (
	"context"

	"wedding-app/domain/model"
)

type UserService interface {
	RegisterUser(ctx context.Context, params model.RegisterUserParams) (*model.User, error)
	LoginUser(ctx context.Context, email, password string) (*model.User, error)
	// RADEK: Tady je super příklad pro ten type alias
	// kdyby si udělal type alias pro ID: type UserID string
	// A potom přejmenoval funkci: Get(ctx context.Context, id UserID) (*model.User, error)
	// Tak máš kratší název funkce a i tak každý kdo příjde ví, že pro Get je potřeba ID uživatele.
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	CreateGuest(ctx context.Context, params model.CreateGuestParams) (*model.User, error)
}
