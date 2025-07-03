package config

import (
	"fmt"
	"go-api/models"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	db, err := gorm.Open(mysql.Open(os.Getenv("DB_URL")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Disable logs
	})

	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}

	err = db.AutoMigrate(
		&models.Permission{},
		&models.Role{},
		&models.User{},
	)

	if err != nil {
		panic(fmt.Sprintf("Auto Migrate Failed: %v", err))
	}

	DB = db
}
