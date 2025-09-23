// RADEK: tohle asi není úplně špatné, ale já bych radši měl domain errory, které se umí samy převést na HTTP errory nebo mít
// konverzi v nějakém middleware, když teď budeš dělat třeba grpc, tak by se to podle mě dělalo líp.
package restapi

import (
	"fmt"
	"net/http"
)

type APIError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
	Err        error  `json:"-"` // The internal error for logging
}

func (e *APIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("status %d: %s, internal error: %v", e.StatusCode, e.Message, e.Err)
	}
	return fmt.Sprintf("status %d: %s", e.StatusCode, e.Message)
}

// NewAPIError creates a new APIError. The 'err' parameter is optional.
func NewAPIError(statusCode int, message string, err error) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
		Err:        err,
	}
}

// NewInternalAPIError The 'err' parameter is optional. RADEK: všechny poznámky musí začínat takto.
func NewInternalAPIError(err error) *APIError {
	return &APIError{
		StatusCode: http.StatusInternalServerError,
		Message:    "internal server error",
		Err:        err,
	}
}

func NewInvalidRequestPayloadAPIError(err error) *APIError {
	return &APIError{
		StatusCode: http.StatusBadRequest,
		Message:    "invalid request payload",
		Err:        err,
	}
}
