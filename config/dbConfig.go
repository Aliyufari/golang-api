package config

import (
	"fmt"
	"gorm.io/gorm"
  	"gorm.io/driver/mysql"
)

func Connect()  {
	const dsn string = "root:@tcp(127.0.0.1:3306)/go_api?charset=utf8mb4&parseTime=True&loc=Local"

	_, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %v", err))
	}
}