package restapi

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	IconUrl  string `json:"iconurl" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type CreateGuestRequest struct {
	Username string `json:"username" binding:"required"`
	IconUrl  string `json:"iconurl" binding:"required"`
	QuizID   string `json:"quizID" binding:"required,uuid"`
}

type CreateQuizRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateQuestionRequest struct {
	Text   string `json:"text" binding:"required"`
	QuizID string `json:"quiz_id" binding:"required,uuid"`
}

type CreateAnswerRequest struct {
	Text       string `json:"text" binding:"required"`
	QuestionID string `json:"question_id" binding:"required,uuid"`
	IsCorrect  bool   `json:"iscorrect"`
}

type SubmitAnswerRequest struct {
	QuestionID string `json:"question_id" binding:"required,uuid"`
	AnswerID   string `json:"answer_id" binding:"required,uuid"`
}
