package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"lotsoflovecindy/m/v2/gcs"
	"lotsoflovecindy/m/v2/model"
)

func main() {
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"}, // Allow necessary methods, including OPTIONS for preflight
		AllowedHeaders: []string{"Content-Type"},
	})

	// Pull credentials from environment variables from docker compose
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Format connection string
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error opening DB: %v", err)
	}

	fmt.Println("Successfully connected to the database!")

	err = db.AutoMigrate(&model.Post{})
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
	fmt.Println("Successfully migrated the database!")

	http.Handle("/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/upload", uploadHandler)

	http.HandleFunc("/list-files", func(w http.ResponseWriter, r *http.Request) {
		err := gcs.RetrieveAllFilesFromGCS(w)
		if err != nil {
			log.Printf("Error retrieving files: %v", err)
		}
	})

	fmt.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler.Handler(http.DefaultServeMux)))
}
