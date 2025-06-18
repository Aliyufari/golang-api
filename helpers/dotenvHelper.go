package helpers

import (
	"log"
	"github.com/joho/godotenv"
)

func LoadEnv() {
	if error := godotenv.Load(); error != nil {
		log.Fatal("Failed to load env")
	}
}