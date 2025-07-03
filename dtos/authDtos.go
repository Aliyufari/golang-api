package dtos

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Dob      string `json:"dob,omitempty" validate:"required,datetime=2006-01-02"`
	Gender   string `json:"gender,omitempty"  validate:"omitempty,oneof=male female"`
	Password string `json:"password" validate:"required,min=6"`
}

func (r *CreateUserRequest) Normalize() {
	r.Gender = strings.ToLower(r.Gender)
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateProfileRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Dob      string `json:"dob,omitempty" validate:"required,datetime=2006-01-02"`
	Gender   string `json:"gender,omitempty"  validate:"omitempty,oneof=male female"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Avatar    string    `json:"avatar"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Dob       string    `json:"dob"`
	Gender    string    `json:"gender"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
