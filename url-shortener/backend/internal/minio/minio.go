package minio

import (
	"context"
	"fmt"
	"io"
	"time"

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

// Save Data into the Mini IO wrapper
// Input:
// - ctx: context.Context
// - key: string (file key/name)
// - r: io.Reader (data to be saved)
// - size: int64 (size of the data)
// - opts: minio.PutObjectOptions (additional options)
// Output:
// - error: any error that occurred during the save operation

func (o *ObjectStore) PutChunks(ctx context.Context, baseKey string, chunks []io.Reader, opts minio.PutObjectOptions) ([]string, error) {
	var urls []string

	for i, r := range chunks {
		// Create a unique key for each chunk: "file.zip.part1", "file.zip.part2", etc.
		chunkKey := fmt.Sprintf("%s.part%d", baseKey, i)

		_, err := o.client.PutObject(ctx, o.bucketName, chunkKey, r, -1, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to upload chunk %d: %w", i, err)
		}

		// Generate URL for this specific chunk
		u, _ := o.client.PresignedGetObject(ctx, o.bucketName, chunkKey, time.Hour, nil)
		urls = append(urls, u.String())
	}

	return urls, nil
}
