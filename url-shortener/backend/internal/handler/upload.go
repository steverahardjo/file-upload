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

func UploadFile(logger *log.Logger, store *minio.ObjectStore, reader io.Reader, db *db.Database, file string, limit int) {
	if

}
