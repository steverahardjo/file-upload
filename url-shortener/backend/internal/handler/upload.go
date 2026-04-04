package handler

import (
	"log"
	"net/http"

	db "github.com/steverahardjo/url-shortener/internal/db"
	minio "github.com/steverahardjo/url-shortener/internal/minio"
)

func HandleUpload(logger *log.Logger, db *db.Database, store *minio.ObjectStore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
