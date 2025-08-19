package restapi

import (
	"errors"
	"net/http"
	"wedding-app/domain/service"

	"github.com/gin-gonic/gin"
)

type QuizHandler struct {
	quizService    service.QuizService
	jwtService     service.JWTService
	sessionService service.SessionService
}

func NewQuizHandler(qs service.QuizService, js service.JWTService) *QuizHandler {
	return &QuizHandler{quizService: qs, jwtService: js} //TODO: add session service
}

// createQuiz godoc
//
//	@Summary		Register a new quiz
//	@Description	Create a quiz with name
//	@Tags			quiz
//
// @Security CookieAuth
//
//	@Accept			json
//	@Produce		json
//	@Param			quiz	body		CreateQuizRequest	true	"Quiz info"
//	@Success		201		{object}	model.Quiz
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/api/create-quiz [post]
func (h *QuizHandler) createQuiz(c *gin.Context) {
	var req CreateQuizRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(NewInvalidRequestPayloadAPIError(err))
		return
	}

	quiz, err := h.quizService.CreateQuiz(c, req.Name)
	if err != nil {
		c.Error(NewInternalAPIError(err))
		return
	}

	c.JSON(http.StatusCreated, quiz)
}

// joinQuiz godoc
//
// @Summary Get a quiz by Invite Code
// @Description Retrieve a single quiz by Invite Code
// @Tags quiz
// @Produce json
// @Param invite query string true "Quiz Invite Code"
// @Success 200 {object} model.Quiz
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/join-quiz [get]
func (h *QuizHandler) joinQuiz(c *gin.Context) {
	inviteCode := c.Query("invite")

	if inviteCode == "" {
		c.Error(NewAPIError(http.StatusBadRequest, "invite query parameter is required", nil))
		return
	}

	quiz, err := h.quizService.GetQuizByInviteCode(c, inviteCode)
	if err != nil {
		c.Error(NewAPIError(http.StatusNotFound, "failed to find quiz!", err))
		return
	}

	authenticated := false
	token, err := c.Cookie(CookieAccessTokenName) // Trying get acces token
	if err == nil {                               // I did find it
		_, verifyErr := h.jwtService.Verify(token) // Trying verify access token
		authenticated = (verifyErr == nil)         // Without error then I know it's authenticated
	}

	// TODO: work differently with response, use unauthorized
	c.JSON(http.StatusOK, gin.H{
		"quiz":          quiz,
		"authenticated": authenticated,
	})
}

// getQuiz godoc
// @Summary Get a quiz by ID
// @Description Retrieve a single quiz by its ID
// @Tags quiz
// @Security CookieAuth
// @Produce json
// @Param id path string true "Quiz ID"
// @Success 200 {object} model.Quiz
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/quiz/{id} [get]
func (h *QuizHandler) getQuiz(c *gin.Context) {
	id := c.Param("id")

	quiz, err := h.quizService.GetQuizByID(c, id)
	if err != nil {
		c.Error(NewAPIError(http.StatusNotFound, "failed to find quiz!", err))
		return
	}

	c.JSON(http.StatusOK, quiz)
}

func (h *QuizHandler) startQuiz(c *gin.Context) {
	quizID := c.Param("quizId")

	userID, err := GetUserIDForQuizFromContext(c, quizID)
	if err != nil {
		if errors.Is(err, ErrUserIsNotAuthorizedForQuizInContext) {
			c.Error(NewAPIError(http.StatusUnauthorized, "unauthorized for quiz!", err))
			return
		}
		c.Error(NewInternalAPIError(err))
		return
	}

	session, err := h.sessionService.StartQuiz(c, userID.String(), quizID)
	if err != nil {
		c.Error(NewInternalAPIError(err))
		return
	}

	question, err := h.sessionService.GetCurrentQuestion(c, session.ID.String())
	if err != nil {
		c.Error(NewInternalAPIError(err))
		return
	}

	c.JSON(http.StatusOK, question)
}
