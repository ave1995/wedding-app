package restapi

import (
	"errors"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type GinHandlers struct {
	User           *UserHandler
	Basic          *BasicHandler
	Quiz           *QuizHandler
	AuthMiddleware gin.HandlerFunc
}

func NewGinHandlers(user *UserHandler, basic *BasicHandler, quiz *QuizHandler, authMiddleware gin.HandlerFunc) (*GinHandlers, error) {
	if user == nil {
		return nil, errors.New("user handler must not be nil")
	}
	if basic == nil {
		return nil, errors.New("basic handler must not be nil")
	}
	if quiz == nil {
		return nil, errors.New("quiz handler must not be nil")
	}
	if authMiddleware == nil {
		return nil, errors.New("auth middleware must not be nil")
	}

	return &GinHandlers{
		User:           user,
		Basic:          basic,
		Quiz:           quiz,
		AuthMiddleware: authMiddleware,
	}, nil
}

func (h *GinHandlers) RegisterAll(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	basic := router.Group("/")
	h.Basic.Register(basic)

	auth := router.Group("/auth")
	h.User.Register(auth)
	h.Quiz.RegisterAnonymous(auth)

	api := router.Group("/api")
	api.Use(h.AuthMiddleware)
	h.Quiz.Register(api)
}
