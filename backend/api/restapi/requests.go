package restapi

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type CreateQuizRequest struct {
	Name string `json:"name" binding:"required"`
}
