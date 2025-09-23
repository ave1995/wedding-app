package apperrors

import "errors"

// RADEK: errory měj na urovni domain package, tenhle package je tu zbytečný

var (
	ErrNotFound           = errors.New("not found")
	ErrSessionCompleted   = errors.New("session already completed")
	ErrAnswerNotSubmitted = errors.New("no answer submitted")
)
