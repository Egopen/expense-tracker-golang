package initializer

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(postgres.Open(os.Getenv("DB_PATH")), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connct to database")
		panic("Failed to connct to database")
	}
}
