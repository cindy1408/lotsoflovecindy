package main

import (
	"fmt"
	"log"
	"net/http"

	"lotsoflovecindy/m/v2/gcs"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request at /upload")

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from request", http.StatusBadRequest)
		return
	}
	defer file.Close() //nolint:errcheck

	log.Printf("File received: %s", header.Filename)

	err = gcs.UploadFileToGCS(w, file, header.Filename)
	if err != nil {
		return
	}

	if _, err := fmt.Fprintf(w, "File uploaded successfully: %s", header.Filename); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}
