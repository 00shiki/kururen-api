package users

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}
