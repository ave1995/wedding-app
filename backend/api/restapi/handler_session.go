package restapi

import (
	"errors"
	"net/http"
	"wedding-app/domain/service"

	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
	sessionService service.SessionService
}

func NewSessionHandler(ss service.SessionService) *SessionHandler {
	return &SessionHandler{sessionService: ss}
}

func (h *SessionHandler) startSession(c *gin.Context) {
	quizID := c.Param("quiz_id")

	userID, err := GetUserIDForQuizFromContext(c, quizID)
	if err != nil {
		if errors.Is(err, ErrUserIsNotAuthorizedForQuizInContext) {
			c.Error(NewAPIError(http.StatusUnauthorized, "unauthorized for quiz!", err))
			return
		}
		c.Error(NewInternalAPIError(err))
		return
	}

	session, err := h.sessionService.StartSession(c, userID.String(), quizID)
	if err != nil {
		c.Error(NewInternalAPIError(err))
		return
	}

	question, err := h.sessionService.GetCurrentQuestion(c, session.ID.String())
	if err != nil {
		c.Error(NewInternalAPIError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id": session.ID,
		"question":   question,
	})
}

func (h *SessionHandler) submitAnswer(c *gin.Context) {
	sessionID := c.Param("session_id")

	var req SubmitAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(NewAPIError(http.StatusBadRequest, "invalid request body", err))
		return
	}

	isCompleted, err := h.sessionService.SubmitAnswer(c, sessionID, req.QuestionID, req.AnswerIDs)
	if err != nil {
		c.Error(NewInternalAPIError(err))
		return
	}

	if isCompleted {
		result, err := h.sessionService.GetResult(c, sessionID)
		if err != nil {
			c.Error(NewInternalAPIError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"completed": true,
			"result": gin.H{
				"score":      result.Score,
				"total":      result.Total,
				"percentage": result.Percentage,
			},
		})
		return
	}

	nextQuestion, err := h.sessionService.GetCurrentQuestion(c, sessionID)
	if err != nil {
		c.Error(NewInternalAPIError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"completed": false,
		"question":  nextQuestion,
	})
}

func (h *SessionHandler) getCurrentQuestion(c *gin.Context) {
	sessionID := c.Param("session_id")

	// načtení session a kontrola uživatele
	session, err := h.sessionService.GetSessionByID(c, sessionID)
	if err != nil {

		c.Error(NewInternalAPIError(err))
		return
	}
	if session.IsCompleted {
		c.JSON(http.StatusOK, gin.H{
			"completed": true,
			"question":  nil,
		})
		return
	}

	_, err = GetUserIDForQuizFromContext(c, session.QuizID.String())
	if err != nil {
		if errors.Is(err, ErrUserIsNotAuthorizedForQuizInContext) {
			c.Error(NewAPIError(http.StatusUnauthorized, "unauthorized for session!", err))
			return
		}
		c.Error(NewInternalAPIError(err))
		return
	}

	question, err := h.sessionService.GetCurrentQuestion(c, sessionID)
	if err != nil {
		c.Error(NewInternalAPIError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"completed": false,
		"question":  question,
	})
}

func (h *SessionHandler) getResult(c *gin.Context) {
	sessionID := c.Param("session_id")

	session, err := h.sessionService.GetSessionByID(c, sessionID)
	if err != nil {
		c.Error(NewInternalAPIError(err))
		return
	}

	_, err = GetUserIDForQuizFromContext(c, session.QuizID.String())
	if err != nil {
		if errors.Is(err, ErrUserIsNotAuthorizedForQuizInContext) {
			c.Error(NewAPIError(http.StatusUnauthorized, "unauthorized for session!", err))
			return
		}
		c.Error(NewInternalAPIError(err))
		return
	}

	result, err := h.sessionService.GetResult(c, sessionID)
	if err != nil {
		c.Error(NewInternalAPIError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"completed": true,
		"result": gin.H{
			"score":      result.Score,
			"total":      result.Total,
			"percentage": result.Percentage,
		},
	})
}
