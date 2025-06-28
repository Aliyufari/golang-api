package requests

import "strings"

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
