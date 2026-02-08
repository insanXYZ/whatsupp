package dto

type User struct {
	ID    int    `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Image string `json:"image,omitempty"`
	Bio   string `json:"bio,omitempty"`
}

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
