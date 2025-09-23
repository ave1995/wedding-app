package service

import (
	"context"

	"wedding-app/domain/model"
)

// RADEK: Zase tip, začíná být celkem nežádoucí použávat slovíčko GET a začínají se používat fráze jako:
// Load, Fetch, Save
// Plus pokud ta funkce nevytváří nebo neloaduje něco speciálního co by nemuselo být jasné tak se dá předpokládat, že služba AnswerService bude pracovat s answers i podle return parametrů takže:
// CreateAnswer => Create
// GetAnswerByID => Fetch
// GetAnswersByQuestionID => FetchByQuestionID
// min kodu líp kodu :D
// Tohle bude platit pro všechno at to nepíšu 10x
type AnswerService interface {
	CreateAnswer(ctx context.Context, text, questionID string, isCorrect bool) (*model.Answer, error)
	GetAnswerByID(ctx context.Context, id string) (*model.Answer, error)
	GetAnswersByQuestionID(ctx context.Context, questionID string) ([]*model.Answer, error)
}
