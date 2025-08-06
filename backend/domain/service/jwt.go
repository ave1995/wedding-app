package service

import "wedding-app/domain/model"

type JWTService interface {
	Generate(user *model.User) (*model.AccessToken, error)
	Verify(tokenString string) (*model.AccessToken, error)
}
