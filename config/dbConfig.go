package config

import (
	"fmt"
	"go-api/models"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	db, err := gorm.Open(mysql.Open(os.Getenv("DB_URL")), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}

	DB = db

	db.AutoMigrate(&models.User{})
}
