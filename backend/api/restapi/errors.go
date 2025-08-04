package restapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// TODO: add switch somehow
func respondWithError(c *gin.Context, httpStatus int, err error, clientMessage string) {
	c.Error(err)

	errMsg := clientMessage
	if errMsg == "" {
		errMsg = "An unexpected error occurred."
	}

	c.JSON(httpStatus, errorResponse{
		Message: errMsg,
	})
}

const InternalServerError = "Internal Server Error"

// respondWithInternalError is a convenience function for 500 errors.
func respondWithInternalError(c *gin.Context, err error) {
	respondWithError(c, http.StatusInternalServerError, err, InternalServerError)
}

// respondWithBadRequest is a convenience function for 400 errors.
func respondWithBadRequest(c *gin.Context, err error, clientMessage string) {
	respondWithError(c, http.StatusBadRequest, err, clientMessage)
}
