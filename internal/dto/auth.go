package dto

// RegisterRequest 注册请求体
type RegisterRequest struct {
	Username string `json:"username" binding:"required" example:"user"`
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"123456"`
}
