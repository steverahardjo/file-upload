package handler

import (
	"log"
	"net/http"

	storage "github.com/url-shortener/backend/internal/db"
)

func HandleUpload(logger *log.Logger, db *storage.Database, store *storage.ObjectStore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
