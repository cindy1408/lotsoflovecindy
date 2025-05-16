package postgres

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"lotsoflovecindy/m/v2/models"
	"os"
	"time"
)

func Connection() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	fmt.Println("DB ENV:", host, port, user, password, dbname)

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	var db *gorm.DB
	var err error

	// Retry for up to 10 seconds
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			fmt.Println("Successfully connected to the database!")
			break
		}
		fmt.Printf("Retrying DB connection (%d/10): %v\n", i+1, err)
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		log.Fatalf("Error opening DB: %v", err)
	}

	err = db.AutoMigrate(&models.Post{})
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
	fmt.Println("Successfully migrated the database!")

	return db, nil
}
