package gcs

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"strings"
)

func DeleteFileFromGCS(objectName string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)
	const prefix = "https://storage.googleapis.com/lotsoflovecindy/"
	filename := strings.TrimPrefix(objectName, prefix)

	object := bucket.Object(filename)

	fmt.Println("Deleting object:", objectName)

	err = object.Delete(ctx) // capture delete error
	if err != nil {
		return err // return the actual delete error
	}

	return nil
}
