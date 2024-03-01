package util

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/sndzhng/gin-template/internal/config"
	"github.com/sndzhng/gin-template/internal/datastore"
	"google.golang.org/api/iterator"
)

func UploadCloudStorageImageBase64(imageBase64 string, imageFolder, imageName string) error {
	if strings.Contains(imageBase64, ",") {
		imageBase64 = strings.Split(imageBase64, ",")[1]
	}
	imageByte, err := base64.StdEncoding.DecodeString(imageBase64)
	if err != nil {
		return err
	}

	imageBuffer := bytes.NewBuffer(imageByte)
	contentType := http.DetectContentType(imageByte)
	imageImage := new(image.Image)

	switch contentType {
	case "image/jpeg":
		*imageImage, err = jpeg.Decode(imageBuffer)
		if err != nil {
			return err
		}
	case "image/png":
		*imageImage, err = png.Decode(imageBuffer)
		if err != nil {
			return err
		}
	default:
		return errors.New("incorrect image type")
	}

	context, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	objectName := fmt.Sprintf("%s/%s.%s", imageFolder, imageName, strings.Split(contentType, "/")[1])
	writer := datastore.CloudStorage.
		Bucket(config.Datastore.CloudStorage.BucketName).
		Object(objectName).
		NewWriter(context)
	err = jpeg.Encode(writer, *imageImage, nil)
	if err != nil {
		return err
	}
	defer writer.Close()

	return nil
}

func GetCloudStorageImageURL(fileFolder, fileName string) (*string, error) {
	query := &storage.Query{
		Prefix: fmt.Sprintf("%s/%s.", fileFolder, fileName),
	}
	objectInterator := datastore.CloudStorage.
		Bucket(config.Datastore.CloudStorage.BucketName).
		Objects(context.Background(), query)
	objectAttrs, err := objectInterator.Next()
	if err != nil {
		switch err {
		case iterator.Done:
			return nil, nil
		default:
			return nil, err
		}
	}

	signedURLOption := &storage.SignedURLOptions{
		GoogleAccessID: datastore.CloudStorageJWTConfig.Email,
		PrivateKey:     datastore.CloudStorageJWTConfig.PrivateKey,
		Method:         http.MethodGet,
		Expires:        time.Now().Add(24 * time.Hour),
	}

	url, err := datastore.CloudStorage.
		Bucket(config.Datastore.CloudStorage.BucketName).
		SignedURL(objectAttrs.Name, signedURLOption)
	if err != nil {
		return nil, err
	}

	return &url, nil
}
