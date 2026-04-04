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
	ShortCode string
	FileType  string
	Size      int
	CreatedAt time.Time
}

type Chunk struct {
	ObjectKey string
	Index     int
}

func NewDatabase(cfg *config.DBConfig) (*Database, error) {
	db, err := sql.Open("postgres", cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &Database{DB: db}, nil
}

func (d *Database) InitSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS files (
		short_code TEXT PRIMARY KEY,
		file_type TEXT NOT NULL,
		bytes INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS file_chunks (
		id SERIAL PRIMARY KEY,
		short_code TEXT NOT NULL,
		chunk_index INTEGER NOT NULL,
		object_key TEXT NOT NULL
	);
	`
	_, err := d.DB.Exec(query)
	return err
}

func (d *Database) CreateFile(file *File) error {
	_, err := d.DB.Exec(`
		INSERT INTO files (short_code, file_type, bytes)
		VALUES ($1, $2, $3)
	`, file.ShortCode, file.FileType, file.Size)
	return err
}

func (d *Database) AddChunk(shortCode string, index int, objectKey string) error {
	_, err := d.DB.Exec(`
		INSERT INTO file_chunks (short_code, chunk_index, object_key)
		VALUES ($1, $2, $3)
	`, shortCode, index, objectKey)
	return err
}

func (d *Database) GetFile(shortCode string) (*File, error) {
	row := d.DB.QueryRow(`
		SELECT short_code, file_type, bytes, created_at
		FROM files
		WHERE short_code = $1
	`, shortCode)

	var f File
	err := row.Scan(&f.ShortCode, &f.FileType, &f.Size, &f.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &f, err
}

func (d *Database) GetChunks(shortCode string) ([]Chunk, error) {
	rows, err := d.DB.Query(`
		SELECT object_key, chunk_index
		FROM file_chunks
		WHERE short_code = $1
		ORDER BY chunk_index ASC
	`, shortCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chunks []Chunk

	for rows.Next() {
		var c Chunk
		if err := rows.Scan(&c.ObjectKey, &c.Index); err != nil {
			return nil, err
		}
		chunks = append(chunks, c)
	}

	return chunks, rows.Err()
}
