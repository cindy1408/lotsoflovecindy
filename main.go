package main

import (
	"fmt"
	"log"
	"lotsoflovecindy/m/v2/gcs"
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

	err = gcs.UploadFileToGCS(w, file, header.Filename)
	if err != nil {
		return
	}

	fmt.Fprintf(w, "File uploaded successfully: %s", header.Filename)
}
