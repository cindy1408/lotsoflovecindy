package postgres

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"lotsoflovecindy/m/v2/models"
	"os"
)

// Connection connects to server
func Connection() (*gorm.DB, error) {
	// Use POSTGRES_ environment variables directly with defaults
	host := getEnvWithDefault("POSTGRES_HOST", "localhost")
	port := getEnvWithDefault("POSTGRES_PORT", "5432")
	user := getEnvWithDefault("POSTGRES_USER", "user")
	password := getEnvWithDefault("POSTGRES_PASSWORD", "password")
	dbname := getEnvWithDefault("POSTGRES_DB", "lotsoflovecindy")

	fmt.Println("DB ENV:", host, port, user, password, dbname)

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	var db *gorm.DB
	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err == nil {
		fmt.Println("Successfully connected to the database!")
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

// Helper function to get environment variable with fallback
func getEnvWithDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
