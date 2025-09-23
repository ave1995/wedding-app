package model

import (
	"encoding/json"

	"github.com/google/uuid"
)

// Session which User use to complete Quiz
type Session struct {
	ID            uuid.UUID `json:"session_id"`
	UserID        uuid.UUID `json:"-"`
	QuizID        uuid.UUID `json:"-"`
	CurrentQIndex int       `json:"-"`
	TotalQCount   int       `json:"-"`
	IsCompleted   bool      `json:"is_completed"`
}

// RADEK: když pominu to, že tahle funkce se nepoužívá tak její přidaná hodnota je 0, jelikož je to marshal structu což víceměně umí každý struct
func (s Session) HttpResponse() ([]byte, error) {
	result, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return result, nil
}
