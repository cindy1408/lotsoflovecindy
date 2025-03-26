package main

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/upload", uploadHandler)

	fmt.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Printf("Failed to copy file to GCS: %v", err)
		http.Error(w, "Failed to create storage client", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	bucketName := "lotsoflovecindy"
	objectName := header.Filename
	bucket := client.Bucket(bucketName)
	object := bucket.Object(objectName)

	writer := object.NewWriter(ctx)
	defer writer.Close()
	if err := writer.Close(); err != nil {
		log.Printf("Failed to close GCS writer: %v", err)
		http.Error(w, fmt.Sprintf("Failed to finalize upload to GCS: %v", err), http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(writer, file)
	if err != nil {
		http.Error(w, "Failed to upload file to GCS", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully: %s", objectName)
}
