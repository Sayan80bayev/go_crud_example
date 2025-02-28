package s3

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go_crud_example/internal/config"
	"log"
	"net/http"
)

func Init(cfg *config.Config) *minio.Client {
	endpoint := "localhost:9000"
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to minio")
	return minioClient
}

func UploadFile(c *gin.Context, cfg *config.Config, minioClient *minio.Client) (string, error) {
	// Get file from request
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		// If no file is uploaded, return an empty string without an error
		if errors.Is(err, http.ErrMissingFile) {
			return "", nil
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return "", errors.New("invalid file")
	}
	defer file.Close()

	objectName := header.Filename
	contentType := header.Header.Get("Content-Type")

	bucketName := cfg.MinioBucket
	_, err = minioClient.PutObject(
		context.Background(),
		bucketName,
		objectName,
		file,
		header.Size,
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return "", errors.New("failed to upload file")
	}

	// Generate file URL
	fileURL := fmt.Sprintf("http://localhost:9000/%s/%s", bucketName, objectName)

	return fileURL, nil
}
