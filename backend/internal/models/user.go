package models

import (
	"time"
)

type UserRole string

const (
	RoleViewer UserRole = "viewer"
	RoleEditor UserRole = "editor"
	RoleAdmin  UserRole = "admin"
)

type User struct {
	ID        string    `db:"id"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Role      UserRole  `db:"role"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type CreateUserRequest struct {
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=8"`
	Role     UserRole `json:"role" binding:"required,oneof=viewer editor admin"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID    string   `json:"id"`
		Email string   `json:"email"`
		Role  UserRole `json:"role"`
	} `json:"user"`
}

type UpdateUserRequest struct {
	Email string   `json:"email" binding:"email"`
	Role  UserRole `json:"role" binding:"oneof=viewer editor admin"`
}
