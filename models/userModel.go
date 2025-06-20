package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Avatar   string `json:"avatar,omitempty" validate:"avatar_format=jpg,png,gif"`
	Name     string `json:"name" validate:"required,min=3,max=32"`
	Email    string `json:"email" gorm:"unique" validate:"required,email,unique"`
	Dob      string `json:"dob,omitempty" validate:"required,date"`
	Gender   string `json:"gender,omitempty"  validate:"omitempty,oneof=male female"`
	Password string `json:"password" validate:"required,min=6"`
}
