package gcs

import (
	"context"
	"fmt"
	"gallery/models"
	"io"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
)

func RetrieveAllPosts(posts []models.Post) ([]models.Post, error) {
	gcsBaseURL := "https://storage.googleapis.com/" + BucketName + "/"

	for i, post := range posts {
		if len(post.ContentURL) > len(gcsBaseURL) && post.ContentURL[:len(gcsBaseURL)] == gcsBaseURL {
			objectName := post.ContentURL[len(gcsBaseURL):]
			signedURL, err := GenerateSignedURL(BucketName, objectName)
			if err != nil {
				log.Printf("Failed to sign URL for object %s: %v", objectName, err)
				continue
			}
			posts[i].ContentURL = signedURL
		}
	}

	return posts, nil
}

func UploadPost(filename, contentType string) (string, error) {
	signedURL, err := GenerateUploadSignedUploadURL(BucketName, filename, contentType)
	if err != nil {
		log.Printf("Failed to generate signed URL: %v", err)
		return "", nil
	}

	return signedURL, nil
}

func DeletePost(decoded string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Println("Error creating storage client:", err)
		return nil
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	o := client.Bucket(BucketName).Object(decoded)

	attrs, err := o.Attrs(ctx)
	if err != nil {
		log.Println("Error getting object attributes:", err)
		return err
	}

	o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})
	if err = o.Delete(ctx); err != nil {
		log.Println("Error deleting GCS object:", err)
		return err
	}

	log.Println("Successfully deleted from GCS")
	return nil
}

func UploadFileToGCS(w http.ResponseWriter, file io.Reader, fileName string) (string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Printf("Failed to create GCS client: %v", err)
		http.Error(w, "Failed to create storage client", http.StatusInternalServerError)
		return "", err
	}
	defer client.Close() //nolint:errcheck

	bucket := client.Bucket(BucketName)
	object := bucket.Object(fileName)
	writer := object.NewWriter(ctx)

	_, err = io.Copy(writer, file)
	if err != nil {
		http.Error(w, "Failed to upload file to GCS", http.StatusInternalServerError)
		return "", err
	}

	if err = writer.Close(); err != nil {
		log.Printf("Failed to close GCS writer: %v", err)
		http.Error(w, fmt.Sprintf("Failed to finalize upload to GCS: %v", err), http.StatusInternalServerError)
		return "", err
	}

	// Return public URL
	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", BucketName, fileName)
	return publicURL, nil
}
