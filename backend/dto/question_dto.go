package dto

import (
	"wedding-app/domain/model"

	"github.com/google/uuid"
)

type QuestionResponse struct {
	SessionID     uuid.UUID       `json:"session_id"`
	Completed     bool            `json:"completed"`
	Question      *model.Question `json:"question"`
	CurrentQIndex int             `json:"currentQIndex"`
	TotalQCount   int             `json:"totalQCount"`
}
