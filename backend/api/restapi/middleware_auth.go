package restapi

import (
	"errors"
	"net/http"
	"wedding-app/domain/model"
	"wedding-app/domain/service"

	"github.com/gin-gonic/gin"
)

const ContextKey = "user"

func AuthMiddleware(jwtService service.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, err := ctx.Cookie("access_token")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		accessToken, err := jwtService.Verify(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Invalid or expired token"})
			return
		}

		ctx.Set(ContextKey, accessToken.UserID)

		ctx.Next()
	}
}

var ErrAccessTokenNotInContext = errors.New("access token not found in context")
var ErrAccessTokenHasInvalidTypeInContext = errors.New("access token in context has invalid type")

func GetUserFromContext(c *gin.Context) (*model.AccessToken, error) {
	val, exists := c.Get(ContextKey)
	if !exists {
		return nil, ErrAccessTokenNotInContext
	}

	accessToken, ok := val.(*model.AccessToken)
	if !ok {
		return nil, ErrAccessTokenHasInvalidTypeInContext
	}

	return accessToken, nil
}
