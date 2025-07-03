package seeders

import (
	"fmt"
	"go-api/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedRolesAndPermissions(db *gorm.DB) error {
	// Define permissions
	permissions := []models.Permission{
		{ID: uuid.New(), Name: "create_user"},
		{ID: uuid.New(), Name: "update_user"},
		{ID: uuid.New(), Name: "delete_user"},
		{ID: uuid.New(), Name: "view_user"},
	}

	// Insert permissions (if not exists)
	for _, p := range permissions {
		var existing models.Permission
		if err := db.Where("name = ?", p.Name).First(&existing).Error; err != nil {
			if err := db.Create(&p).Error; err != nil {
				return fmt.Errorf("failed to seed permission %s: %v", p.Name, err)
			}
		}
	}

	// Fetch permissions from DB
	var allPermissions []models.Permission
	db.Find(&allPermissions)

	// Create admin role and attach all permissions
	adminRole := models.Role{
		ID:          uuid.New(),
		Name:        "admin",
		Permissions: allPermissions,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Check and insert admin role
	var existingAdmin models.Role
	if err := db.Where("name = ?", adminRole.Name).First(&existingAdmin).Error; err != nil {
		if err := db.Create(&adminRole).Error; err != nil {
			return fmt.Errorf("failed to seed admin role: %v", err)
		}
	}

	// Create user role (no permissions)
	userRole := models.Role{
		ID:        uuid.New(),
		Name:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	var existingUser models.Role
	if err := db.Where("name = ?", userRole.Name).First(&existingUser).Error; err != nil {
		if err := db.Create(&userRole).Error; err != nil {
			return fmt.Errorf("failed to seed user role: %v", err)
		}
	}

	fmt.Println("✅ Roles and permissions seeded successfully")
	return nil
}
