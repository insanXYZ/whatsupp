package dto

type Response struct {
	Message string `json:"message,omitempty"`
	Error   any    `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
}
