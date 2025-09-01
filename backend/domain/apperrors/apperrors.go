package apperrors

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrSessionCompleted   = errors.New("session already completed")
	ErrAnswerNotSubmitted = errors.New("no answer submitted")
)
