package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	ID          uuid.UUID    `gorm:"primaryKey;type:char(36)" json:"id"`
	Name        string       `gorm:"type:varchar(100);unique" json:"name"`
	Permissions []Permission `gorm:"many2many:permission_role;" json:"permissions"`
	Users       []User       `gorm:"foreignKey:RoleID" json:"users"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

func (u *Role) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}
