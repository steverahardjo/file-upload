package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	config "github.com/url-shortener/backend/internal/config"
)

type Database struct {
	DB *sql.DB
}

type File struct {
	ID            int       `db:"id"`
	URL           string    `db:"url"`
	ShortCode     string    `db:"short_code"`
	CreatedAt     time.Time `db:"created_at"`
	Size          int       `db:"bytes"`
	FileType      string    `db:"file_type"`
	ChunksGroupID int       `db:"chunks_group_id"`
}

func NewDatabase(cfg *config.DBConfig) (*Database, error) {
	db, err := sql.Open("postgres", cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &Database{DB: db}, nil
}

func (d *Database) InitSchema() error {
	query := `
		CREATE TABLE IF NOT EXISTS files (
			id SERIAL PRIMARY KEY,
			url TEXT NOT NULL,
			short_code TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			bytes INTEGER NOT NULL,
			file_type TEXT NOT NULL,
			chunks_group_id INTEGER NOT NULL
		);
	`
	_, err := d.DB.Exec(query)
	return err
}

func (d *Database) AddFileRecord(file *File) error {
	query := `
		INSERT INTO files (url, short_code, bytes, file_type, chunks_group_id)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := d.DB.Exec(query, file.URL, file.ShortCode, file.Size, file.FileType, file.ChunksGroupID)
	return err
}
