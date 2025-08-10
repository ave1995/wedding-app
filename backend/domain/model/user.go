package model

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID
	Username string
	Email    string
	IconUrl  string
	IsGuest  bool
	QuizID   uuid.UUID
}

type RegisterUserParams struct {
	Username string
	Email    string
	Password string
	IconURL  string
}

type CreateGuestParams struct {
	Username string
	IconUrl  string
	QuizID   string
}
