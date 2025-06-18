package config

import (
	"fmt"
	"os"
	"gorm.io/gorm"
  	"gorm.io/driver/mysql"
)

var DB *gorm.DB

func ConnectDB()  {
	db, err := gorm.Open(mysql.Open(os.Getenv("DB_URL")), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}

	DB = db
}