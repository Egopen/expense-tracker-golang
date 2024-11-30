package initializer

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if _, exists := os.LookupEnv("PORT"); !exists {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}
