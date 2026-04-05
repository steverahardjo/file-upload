package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/steverahardjo/url-shortener/internal/db"
	"github.com/steverahardjo/url-shortener/internal/minio"
)

func (h *Handler) HandleDownload() http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		shortCode := strings.TrimPrefix(r.URL.Path, "/download/")

		h.logging.Printf("download request: %s", shortCode)

		err := DownloadHelper(r.Context(), h.db, h.obj_store, shortCode, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func DownloadHelper(
	ctx context.Context,
	database *db.Database,
	store *minio.ObjectStore,
	shortCode string,
	w http.ResponseWriter,
) error {
	file, err := database.GetFile(shortCode)
	if err != nil {
		return err
	}
	if file == nil {
		return fmt.Errorf("file not found")
	}

	chunks, err := database.GetChunks(shortCode)
	if err != nil {
		return err
	}
	if len(chunks) == 0 {
		return fmt.Errorf("no chunks found")
	}

	// pre-allocate by index — no mutex needed
	type result struct {
		data []byte
		err  error
	}
	results := make([]result, len(chunks))

	var wg sync.WaitGroup
	for i, c := range chunks {
		wg.Add(1)
		go func(idx int, objectKey string) {
			defer wg.Done()

			reader, err := store.GetObject(ctx, objectKey)
			if err != nil {
				results[idx].err = err
				return
			}
			defer reader.Close()

			data, err := io.ReadAll(reader)
			results[idx] = result{data: data, err: err}
		}(i, c.ObjectKey)
	}
	wg.Wait()

	// write headers only after all chunks are verified
	for _, r := range results {
		if r.err != nil {
			return r.err
		}
	}

	w.Header().Set("Content-Type", file.FileType)
	w.Header().Set("Content-Disposition", "attachment; filename="+file.ShortCode)

	// stream to response in order
	for _, r := range results {
		if _, err := w.Write(r.data); err != nil {
			return err
		}
	}
	return nil
}
