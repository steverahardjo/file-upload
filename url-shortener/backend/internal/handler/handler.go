package handler

import (
	"log"

	db "github.com/steverahardjo/url-shortener/internal/db"
	minio "github.com/steverahardjo/url-shortener/internal/minio"
)

type Handler struct {
	logging   *log.Logger
	db        *db.Database
	obj_store *minio.ObjectStore
}
