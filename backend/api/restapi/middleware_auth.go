package restapi

import (
	"errors"
	"net/http"
	"wedding-app/domain/model"
	"wedding-app/domain/service"

	"github.com/gin-gonic/gin"
)

const CookieAccessTokenName = "access_token"

const ContextUserIDKey = "userID"

func AuthMiddleware(jwtService service.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, err := ctx.Cookie(CookieAccessTokenName)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		accessToken, err := jwtService.Verify(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Invalid or expired token"})
			return
		}

		// TODO: I have to know how to control if user is guest
		ctx.Set(ContextUserIDKey, accessToken.UserID)

		ctx.Next()
	}
}

var ErrAccessTokenNotInContext = errors.New("access token not found in context")
var ErrAccessTokenHasInvalidTypeInContext = errors.New("access token in context has invalid type")

func GetUserFromContext(c *gin.Context) (*model.AccessToken, error) {
	val, exists := c.Get(ContextUserIDKey)
	if !exists {
		return nil, ErrAccessTokenNotInContext
	}

	accessToken, ok := val.(*model.AccessToken)
	if !ok {
		return nil, ErrAccessTokenHasInvalidTypeInContext
	}

	return accessToken, nil
}
