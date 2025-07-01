package main

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Role     string `json:"role" validate:"required"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type PingResponse struct {
	Message string `json:"message"`
	Service string `json:"service"`
}

type JWTResponse struct {
	Token string `json:"token"`
}

type BlockUserRequest struct {
	Username string `json:"username" validate:"required"`
}
