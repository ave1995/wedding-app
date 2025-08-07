package restapi

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandler is a middleware that catches errors pushed to the context
func ErrorHandler(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// After the request has been processed, check for errors
		if len(c.Errors) > 0 {
			// Get the last error
			err := c.Errors.Last()

			// Check if it's our custom APIError
			if apiErr, ok := err.Err.(*APIError); ok {
				// Log the internal error for debugging
				// You can replace this with your preferred logging library (e.g., zap, logrus)
				fmt.Println(apiErr.Error())

				// Send the client-facing response
				c.AbortWithStatusJSON(apiErr.StatusCode, apiErr)
			} else {
				// Handle unhandled or unexpected errors
				// Log the raw error
				fmt.Println("Unhandled error:", err.Err)

				// Send a generic 500 error to the client
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "An unexpected error occurred.",
				})
			}
		}
	}
}
