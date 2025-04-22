package initializers

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDb() {

	var err error

	dsn := os.Getenv("DB")

	if dsn == "" {
		log.Fatal("DB not variable")
	}

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Panic("Problem connect DB")
	}

	log.Println("DB connect +")
}
