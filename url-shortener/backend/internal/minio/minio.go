package minio

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	config "github.com/steverahardjo/url-shortener/internal/config"
)

type ObjectStore struct {
	client     *minio.Client
	bucketName string
	size_limit int
}

func NewFromConfig(cfg config.ObjStoreConfig) (*ObjectStore, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.Secure,
	})
	if err != nil {
		return nil, err
	}

	return &ObjectStore{
		client:     client,
		bucketName: cfg.BucketName,
	}, nil
}

func (o *ObjectStore) GetObject(ctx context.Context, key string) (io.ReadCloser, error) {
	obj, err := o.client.GetObject(ctx, o.bucketName, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("get object: %w", err)
	}

	// force existence check
	_, err = obj.Stat()
	if err != nil {
		obj.Close()
		return nil, fmt.Errorf("object not found: %w", err)
	}

	return obj, nil
}

func (o *ObjectStore) PutObject(ctx context.Context, key string, r io.Reader, size int64, opts minio.PutObjectOptions) error {
	if size > int64(o.size_limit) {
		return fmt.Errorf("file size exceeds limit: %d > %d", size, o.size_limit)
	}
	_, err := o.client.PutObject(ctx, o.bucketName, key, r, size, opts)
	return err
}
