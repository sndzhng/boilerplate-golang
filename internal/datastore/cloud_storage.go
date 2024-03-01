package datastore

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"github.com/sndzhng/gin-template/internal/config"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
)

var (
	CloudStorage          *storage.Client
	CloudStorageJWTConfig *jwt.Config
)

func ConnectCloudStorage() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "./cloud-storage-credential.json")

	// description: initial cloud storage client
	err := error(nil)
	CloudStorage, err = storage.NewClient(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// description: check bucket exist
	_, err = CloudStorage.Bucket(config.Datastore.CloudStorage.BucketName).Attrs(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// description: get jwt config
	jsonKey, err := os.ReadFile("./cloud-storage-credential.json")
	if err != nil {
		log.Fatal(err)
	}

	CloudStorageJWTConfig, err = google.JWTConfigFromJSON(jsonKey)
	if err != nil {
		log.Fatal(err)
	}
}
