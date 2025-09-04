package restapi

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type GinHandlers struct {
	User           *UserHandler
	Basic          *BasicHandler
	Quiz           *QuizHandler
	Question       *QuestionHandler
	Answer         *AnswerHandler
	Session        *SessionHandler
	WS             *WSHandler
	AuthMiddleware gin.HandlerFunc
}

func NewGinHandlers(user *UserHandler, basic *BasicHandler, quiz *QuizHandler, question *QuestionHandler, answer *AnswerHandler, session *SessionHandler, ws *WSHandler, authMiddleware gin.HandlerFunc) (*GinHandlers, error) {
	if user == nil {
		return nil, errors.New("user handler must not be nil")
	}
	if basic == nil {
		return nil, errors.New("basic handler must not be nil")
	}
	if quiz == nil {
		return nil, errors.New("quiz handler must not be nil")
	}
	if question == nil {
		return nil, errors.New("question handler must not be nil")
	}
	if answer == nil {
		return nil, errors.New("answer handler must not be nil")
	}
	if session == nil {
		return nil, errors.New("session handler must not be nil")
	}
	if ws == nil {
		return nil, errors.New("ws handler must not be nil")
	}
	if authMiddleware == nil {
		return nil, errors.New("auth middleware must not be nil")
	}

	return &GinHandlers{
		User:           user,
		Basic:          basic,
		Quiz:           quiz,
		Question:       question,
		Answer:         answer,
		Session:        session,
		WS:             ws,
		AuthMiddleware: authMiddleware,
	}, nil
}

func (h *GinHandlers) RegisterAll(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	router.GET("/user-svgs", h.Basic.getUserSvgs)

	auth := router.Group("/auth")
	auth.POST("/register", h.User.registerUser)
	auth.POST("/login", h.User.loginUser)
	auth.POST("/create-guest", h.User.createGuest)
	auth.GET("/join-quiz", h.Quiz.joinQuiz)

	api := router.Group("/api")
	api.Use(h.AuthMiddleware)

	api.GET("/auth-check", Require(RoleUser), func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	api.GET("/ws", Require(RoleUser), h.WS.serveWS)

	// Quiz endpoints
	api.POST("/create-quiz", h.Quiz.createQuiz)
	api.GET("/quiz/:id", h.Quiz.getQuiz)
	// Questions endpoints
	api.POST("/create-question", h.Question.createQuestion)
	api.GET("/questions/:id", h.Question.getQuestionByID)
	api.GET("/questions", h.Question.getQuestionsByQuizID)
	// Answers endpoints
	api.POST("/create-answer", h.Answer.createAnswer)
	api.GET("/answers/:id", h.Answer.getAnswerByID)
	api.GET("/answers", h.Answer.getAnswersByQuestionID)

	quiz := api.Group("/quizzes")
	{
		// vytvoření nové session pro daný quiz
		quiz.POST("/:quiz_id/sessions", h.Session.startSession)
	}

	// Session endpoints
	sessions := api.Group("/sessions")
	{
		// odeslání odpovědi
		sessions.POST("/:session_id/answers", h.Session.submitAnswer)

		// načtení aktuální otázky
		sessions.GET("/:session_id/question", h.Session.getCurrentQuestion)

		// získání výsledku
		sessions.GET("/:session_id/result", h.Session.getResult)
	}
}
