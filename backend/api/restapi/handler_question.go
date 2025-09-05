package restapi

import (
	"net/http"
	"strconv"
	"wedding-app/domain/service"

	"github.com/gin-gonic/gin"
)

type QuestionHandler struct {
	questionService service.QuestionService
}

func NewQuestionHandler(qs service.QuestionService) *QuestionHandler {
	return &QuestionHandler{questionService: qs}
}

// createQuestion godoc
//
//	@Summary		Register a new question
//	@Description	Create a question under a specific quiz
//	@Tags			question
//
// @Security CookieAuth
//
//	@Accept			json
//	@Produce		json
//	@Param			question	body		CreateQuestionRequest	true	"Question info"
//	@Success		201		{object}	model.Question
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/api/create-question [post]
func (h *QuestionHandler) createQuestion(c *gin.Context) {
	var req CreateQuestionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(NewInvalidRequestPayloadAPIError(err))
		return
	}

	question, err := h.questionService.CreateQuestion(c, req.Text, req.QuizID, req.Type, req.PhotoPath)
	if err != nil {
		c.Error(NewInternalAPIError(err))
		return
	}

	c.JSON(http.StatusCreated, question)
}

// getQuestionByID godoc
//
//	@Summary		Get a question by ID
//	@Description	Retrieve a question using its unique ID
//	@Tags			question
//
// @Security CookieAuth
//
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Question ID"
//	@Success		200	{object}	model.Question
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/api/questions/{id} [get]
func (h *QuestionHandler) getQuestionByID(c *gin.Context) {
	id := c.Param("id")

	question, err := h.questionService.GetQuestionByID(c, id)
	if err != nil {
		c.Error(NewInternalAPIError(err))
		return
	}

	c.JSON(http.StatusOK, question)
}

// getQuestionsByQuizID godoc
//
//	@Summary		Get all questions for a quiz
//	@Description	Retrieve all questions belonging to a specific quiz
//	@Tags			question
//
// @Security CookieAuth
//
//	@Accept			json
//	@Produce		json
//	@Param			quiz_id	query		string	true	"Quiz ID"
//	@Success		200		{array}	model.Question
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/api/questions [get]
func (h *QuestionHandler) getQuestionsByQuizID(c *gin.Context) {
	quizID := c.Query("quiz_id")

	questions, err := h.questionService.GetQuestionsByQuizID(c, quizID)
	if err != nil {
		c.Error(NewInternalAPIError(err))
		return
	}

	c.JSON(http.StatusOK, questions)
}

func (h *QuestionHandler) revealQuestionByQuizID(c *gin.Context) {
	quizID := c.Param("id")
	index := c.Query("index")

	if index == "" {
		c.Error(NewAPIError(http.StatusBadRequest, "missing index query parameter", nil))
		return
	}

	num, err := strconv.Atoi(index)
	if err != nil {
		c.Error(NewInternalAPIError(err))
		return
	}

	reveal, err := h.questionService.RevealQuestionByQuizIDAndIndex(c, quizID, num)
	if err != nil {
		c.Error(NewInternalAPIError(err))
		return
	}

	c.JSON(http.StatusOK, reveal)
}
