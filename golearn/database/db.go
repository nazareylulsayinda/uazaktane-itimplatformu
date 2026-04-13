package database

import (
	"log"

	"golearn/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	database, err := gorm.Open(sqlite.Open("golearn.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database!", err)
	}

	DB = database

	// Modellerin veritabanı tablolarını otomatik oluşturma (AutoMigrate)
	err = DB.AutoMigrate(
		&models.User{},
		&models.Course{},
		&models.Lesson{},
		&models.Quiz{},
		&models.Question{},
		&models.Progress{},
		&models.QuizResult{},
	)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Database connection established and migrations ok.")
}
