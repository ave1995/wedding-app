package model

import "github.com/google/uuid"

type Quiz struct {
	ID         uuid.UUID
	Name       string
	InviteCode uuid.UUID
}
