package restapi

import (
	"net/http"
	"wedding-app/domain/service"

	"github.com/gin-gonic/gin"
)

type AnswerHandler struct {
	answerService service.AnswerService
}

func NewAnswerHandler(as service.AnswerService) *AnswerHandler {
	return &AnswerHandler{answerService: as}
}

// createAnswer godoc
//
//	@Summary		Register a new answer
//	@Description	Create an answer for a specific question
//	@Tags			answer
//
// @Security CookieAuth
//
//	@Accept			json
//	@Produce		json
//	@Param			answer	body		CreateAnswerRequest	true	"Answer info"
//	@Success		201		{object}	model.Answer
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/api/create-answer [post]
func (h *AnswerHandler) createAnswer(c *gin.Context) {
	var req CreateAnswerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(NewInvalidRequestPayloadAPIError(err))
		return
	}

	answer, err := h.answerService.CreateAnswer(c, req.Text, req.QuestionID, req.IsCorrect)
	if err != nil {
		c.Error(NewInternalAPIError(err))
		return
	}

	c.JSON(http.StatusCreated, answer)
}

// getAnswerByID godoc
//
//	@Summary		Get an answer by ID
//	@Description	Retrieve an answer using its unique ID
//	@Tags			answer
//
// @Security CookieAuth
//
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Answer ID"
//	@Success		200	{object}	model.Answer
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/api/answers/{id} [get]
func (h *AnswerHandler) getAnswerByID(c *gin.Context) {
	id := c.Param("id")

	answer, err := h.answerService.GetAnswerByID(c, id)
	if err != nil {
		c.Error(NewInternalAPIError(err))
		return
	}

	c.JSON(http.StatusOK, answer)
}

// getAnswersByQuestionID godoc
//
//	@Summary		Get all answers for a question
//	@Description	Retrieve all answers belonging to a specific question
//	@Tags			answer
//
// @Security CookieAuth
//
//	@Accept			json
//	@Produce		json
//	@Param			question_id	query		string	true	"Question ID"
//	@Success		200		{array}	model.Answer
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/api/answers [get]
func (h *AnswerHandler) getAnswersByQuestionID(c *gin.Context) {
	questionID := c.Query("question_id")

	answers, err := h.answerService.GetAnswersByQuestionID(c, questionID)
	if err != nil {
		c.Error(NewInternalAPIError(err))
		return
	}

	c.JSON(http.StatusOK, answers)
}
