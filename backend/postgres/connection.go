package postgres

import (
	"fmt"
	"gallery/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

func Connection() (*gorm.DB, error) {
	host := "localhost"
	port := "5433"
	user := "user"
	password := "password"
	dbname := "lotsoflovecindy"

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
	time.Sleep(1 * time.Second)

	if err != nil {
		log.Fatalf("Error opening DB: %v", err)
	}

	err = db.AutoMigrate(&models.Post{}, &models.User{})
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
	fmt.Println("Successfully migrated the database!")

	return db, nil
}
