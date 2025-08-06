package restapi

import (
	"net/http"
	"wedding-app/domain/service"

	"github.com/gin-gonic/gin"
)

type QuizHandler struct {
	quizService service.QuizService
}

func NewQuizHandler(qs service.QuizService) *QuizHandler {
	return &QuizHandler{quizService: qs}
}

func (h *QuizHandler) Register(router *gin.RouterGroup) {
	router.POST("/create-quiz", h.createQuiz)
	router.GET("/join-quiz", h.joinQuiz)
	router.GET("/get-quiz", h.getQuiz)
}

// registerQuiz godoc
//
//	@Summary		Register a new quiz
//	@Description	Create a quiz with name
//	@Tags			quiz
//	@Security BearerAuth
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
		respondWithBadRequest(c, err, "invalid request payload")
		return
	}

	quiz, err := h.quizService.CreateQuiz(c, req.Name)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	c.JSON(http.StatusCreated, quiz)
}

// GetQuizHandler godoc
// @Summary Get a quiz by Invite Code
// @Description Retrieve a single quiz by Invite Code
// @Tags quiz
// @Security BearerAuth
// @Produce json
// @Param invite path string true "Quiz Invite Code"
// @Success 200 {object} model.Quiz
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/quiz/{invite} [get]
func (h *QuizHandler) joinQuiz(c *gin.Context) {
	inviteCode := c.Param("invite")

	quiz, err := h.quizService.GetQuizByInviteCode(c, inviteCode)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, quiz)
}

// GetQuizHandler godoc
// @Summary Get a quiz by ID
// @Description Retrieve a single quiz by its ID
// @Tags quiz
// @Security BearerAuth
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
		respondWithInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, quiz)
}
