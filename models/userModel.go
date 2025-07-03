package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Avatar    string    `json:"avatar"`
	Name      string    `gorm:"type:varchar(100)" json:"name"`
	Email     string    `gorm:"type:varchar(191);unique" json:"email"`
	Dob       string    `json:"dob"`
	Gender    string    `json:"gender"`
	Password  string    `json:"-"`
	RoleID    uuid.UUID `gorm:"type:uuid" json:"role_id"`
	Role      Role      `gorm:"foreignKey:RoleID" json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}
