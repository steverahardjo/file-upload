package handler

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	db "github.com/steverahardjo/url-shortener/internal/db"
	minio "github.com/steverahardjo/url-shortener/internal/minio"
)

func HandleUpload(logger *log.Logger, db *db.Database, store *minio.ObjectStore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Limit the size of the request (e.g., 10MB)
		r.ParseMultipartForm(10 << 20)

		// 2. Retrieve the file from form data
		file, handler, err := r.FormFile("myFile") // "myFile" is the form field name
		if err != nil {
			logger.Printf("Error retrieving the file: %v", err)
			http.Error(w, "Invalid file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// 3. Prepare MinIO options (Metadata & Content-Type)

		opts := minio.ObjectStore.PutObjectOptions{
			ContentType: handler.Header.Get("Content-Type"),
			UserMetadata: map[string]string{
				"x-amz-meta-original-name": handler.Filename,
			},
		}

		// 4. Call the upload logic
		// In a real app, you'd likely generate a unique filename for 'fileKey'
		fileKey := handler.Filename
		UploadFile(r.Context(), logger, store, file, db, fileKey, opts)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Successfully Uploaded: %s", handler.Filename)
	})
}

func UploadFile(
	ctx context.Context,
	logger *log.Logger,
	store *minio.ObjectStore,
	reader io.Reader,
	db *db.Database,
	file string,
	opts minio.PutObjectOptions, // Added parameter
) {
	logger.Printf("Saving file into MinIO: %s", file)

	// Pass the reader and the options into PutChunks
	urls, err := store.PutChunks(ctx, file, []io.Reader{reader}, opts)
	if err != nil {
		logger.Printf("Failed to save file into MinIO: %s", err)
		return
	}

	// 5. Database Logic
	// Use the DB connection to save the record of the upload
	if len(urls) > 0 {
		logger.Printf("File stored at: %s", urls[0])
		// Example: db.Exec("INSERT INTO uploads (name, url) VALUES (?, ?)", file, urls[0])
	}
}
