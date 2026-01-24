package dto

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=8"`
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,gte=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=8"`
}

type UpdateUserRequest struct {
	Name     string `json:"name" validate:"omitempty,gte=3"`
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"omitempty,gte=8"`
}
