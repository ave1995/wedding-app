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
	basic.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	basic.GET("/user-svgs", h.Basic.getUserSvgs)

	auth := router.Group("/auth")
	auth.POST("/register", h.User.registerUser)
	auth.POST("/login", h.User.loginUser)
	auth.POST("/create-guest", h.User.createGuest)
	auth.GET("/join-quiz", h.Quiz.joinQuiz)

	api := router.Group("/api")
	api.Use(h.AuthMiddleware)
	api.POST("/create-quiz", h.Quiz.createQuiz)
	api.GET("/quiz/:id", h.Quiz.getQuiz)
}
