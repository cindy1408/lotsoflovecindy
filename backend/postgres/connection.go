package postgres

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"lotsoflovecindy/m/v2/models"
)

func Connection() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error opening DB: %v", err)
	}

	fmt.Println("Successfully connected to the database!")

	err = db.AutoMigrate(&models.Post{})
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
	fmt.Println("Successfully migrated the database!")

	return db, nil
}
