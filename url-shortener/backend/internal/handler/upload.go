package handler

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/steverahardjo/url-shortener/internal/db"
)

func (h *Handler) HandleUpload() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, "File too large or invalid form", http.StatusRequestEntityTooLarge)
			return
		}

		// 2. Retrieve the file
		file, header, err := r.FormFile("userFile")
		if err != nil {
			h.logging.Printf("Error retrieving file: %v", err)
			http.Error(w, "Invalid file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		shortCode := strings.ReplaceAll(uuid.New().String()[:8], "-", "")

		contentType := header.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		opts := h.obj_store.GetMetadata(contentType, header.Filename)

		_, err = h.obj_store.PutChunks(r.Context(), shortCode, []io.Reader{file}, opts)
		if err != nil {
			h.logging.Printf("MinIO Upload Error: %v", err)
			http.Error(w, "Failed to upload to storage", http.StatusInternalServerError)
			return
		}

		fileRecord := &db.File{
			ShortCode: shortCode,
			FileType:  contentType,
			Size:      int(header.Size),
		}

		if err := h.db.CreateFile(fileRecord); err != nil {
			h.logging.Printf("DB CreateFile Error: %v", err)
			http.Error(w, "Failed to save file record", http.StatusInternalServerError)
			return
		}

		objectKey := fmt.Sprintf("%s.part0", shortCode)
		if err := h.db.AddChunk(shortCode, 0, objectKey); err != nil {
			h.logging.Printf("DB AddChunk Error: %v", err)
			http.Error(w, "Failed to save chunk record", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Successfully Uploaded!\nShort Code: %s\nOriginal Name: %s", shortCode, header.Filename)
	})
}
