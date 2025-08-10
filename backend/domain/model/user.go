package model

import "github.com/google/uuid"

type User struct {
	ID          uuid.UUID
	Username    string
	Email       string
	IconUrl     string
	IsTemporary bool
}

type RegisterUserParams struct {
	Username string
	Email    string
	Password string
	IconURL  string
}
