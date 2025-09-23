package model

import "github.com/google/uuid"

// RADEK: to je jen takový tip, my v oddinu děláme type aliasy i pro IDčka, protože potom hezky vidíš třeba v parametru funkce, co je od tebe potřeba:
// type QuizID uuid.UUID
// func foo (id QuizID)
// func foo (id uuid.UUID)
// Vidíš ten rozdíl, ta funkce ti sama říká co chce za ID
type Quiz struct {
	ID         uuid.UUID
	Name       string
	InviteCode uuid.UUID
}
