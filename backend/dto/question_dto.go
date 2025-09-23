package dto

import (
	"wedding-app/domain/model"

	"github.com/google/uuid"
)

// RADEK: Celé "dto" je něco co nemělo vzniknout
// SessionService.GetCurrentQuestion musí vracet vždycky model a response už je doména RestAPI podobně jako si měl requesty
// To víceméně platí pro všechny ty dto modely co tu máš.
type QuestionResponse struct {
	SessionID     uuid.UUID       `json:"session_id"`
	Completed     bool            `json:"completed"`
	Question      *model.Question `json:"question"`
	CurrentQIndex int             `json:"currentQIndex"`
	TotalQCount   int             `json:"totalQCount"`
}
