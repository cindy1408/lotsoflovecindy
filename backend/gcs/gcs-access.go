package gcs

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"google.golang.org/api/option"
	"os"
	"time"
)

type serviceAccountJSON struct {
	ClientEmail string `json:"client_email"`
	PrivateKey  string `json:"private_key"`
}

const serviceAccount = "../credentials.json"

func GenerateSignedURL(bucket, object string) (string, error) {
	ctx := context.Background()

	credsJSON, err := os.ReadFile(serviceAccount)
	if err != nil {
		return "", err
	}

	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(credsJSON))
	if err != nil {
		return "", err
	}
	defer client.Close()

	var sa serviceAccountJSON
	if err = json.Unmarshal(credsJSON, &sa); err != nil {
		return "", err
	}

	url, err := storage.SignedURL(bucket, object, &storage.SignedURLOptions{
		GoogleAccessID: sa.ClientEmail,
		PrivateKey:     []byte(sa.PrivateKey),
		Method:         "GET",
		Expires:        time.Now().Add(15 * time.Minute),
	})
	if err != nil {
		return "", err
	}

	return url, nil
}

func GenerateUploadSignedUploadURL(bucket, object, contentType string) (string, error) {
	ctx := context.Background()

	credsJSON, err := os.ReadFile(serviceAccount)
	if err != nil {
		return "", err
	}

	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(credsJSON))
	if err != nil {
		return "", err
	}
	defer client.Close()

	var sa serviceAccountJSON
	if err = json.Unmarshal(credsJSON, &sa); err != nil {
		return "", err
	}

	url, err := storage.SignedURL(bucket, object, &storage.SignedURLOptions{
		GoogleAccessID: sa.ClientEmail,
		PrivateKey:     []byte(sa.PrivateKey),
		Method:         "PUT",
		Expires:        time.Now().Add(15 * time.Minute),
		ContentType:    contentType,
	})
	if err != nil {
		return "", err
	}

	return url, nil
}
