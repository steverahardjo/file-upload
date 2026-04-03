package db

import (
	"context"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/steverahardjo/url-shortener/internal/config"
	config "github.com/url-shortener/backend/internal/config"
)

type MinioClient struct {
	client     *minio.Client
	bucketName string
}

func NewMinioClient(ctx context.Context, cfg config.ObjStoreConfig) *MinioClient {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.Secure,
	})
	if err != nil {
		log.Fatal(err)
	}

	m := &MinioClient{
		client:     client,
		bucketName: cfg.BucketName,
	}

	// ensure bucket exists during init
	m.ensureBucket(ctx)

	return m
}

func (c *MinioClient) Add(ctx context.Context, f *os.File, contentType string) (string, error) {
	defer f.Close()

	objectID := uuid.New().String()

	_, err := c.client.FPutObject(
		ctx,
		c.bucketName,
		objectID,
		f.Name(),
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return "", err
	}

	return objectID, nil
}

func (c *MinioClient) Get(ctx context.Context, objectID string) (*os.File, error) {

}
