# URL Shortener with Chunked File Storage (Go, MinIO, PostgreSQL)

## Overview

This project is a backend service written in Go that provides file download via short codes. Files are stored as chunks in object storage (MinIO) and reconstructed on demand using metadata stored in PostgreSQL.

The system is designed with a clear separation of responsibilities:

* HTTP layer handles requests and responses
* Database layer manages metadata and structure
* Object storage handles binary data

---

## Architecture

```
Client
  ↓
HTTP Handler (/download/{shortcode})
  ↓
PostgreSQL (metadata, chunk ordering)
  ↓
MinIO (object storage)
  ↓
Streamed HTTP response
```

---

## Project Structure

```
/cmd/api/main.go           Entry point
/internal/config           Configuration loading
/internal/db               Database access layer
/internal/minio            Object storage abstraction
/internal/handler          HTTP handlers
```

---

## Configuration

The service is configured through environment variables:

```
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minio
MINIO_SECRET_KEY=minio123
MINIO_BUCKET=files
MINIO_SSL=false

DB_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable
```

---

## Database Schema

### files

Stores file-level metadata.

```sql
CREATE TABLE files (
    short_code TEXT PRIMARY KEY,
    file_type TEXT NOT NULL,
    bytes INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### file_chunks

Stores chunk-level information for reconstruction.

```sql
CREATE TABLE file_chunks (
    id SERIAL PRIMARY KEY,
    short_code TEXT NOT NULL,
    chunk_index INTEGER NOT NULL,
    object_key TEXT NOT NULL
);
```

---

## Running the Service

Start the server:

```
go run cmd/api/main.go
```

The service listens on:

```
http://localhost:8080
```

---

## API

### Download File

```
GET /download/{shortcode}
```

#### Behavior

1. Extract `short_code` from the request path
2. Query file metadata from PostgreSQL
3. Query chunk list ordered by `chunk_index`
4. Stream each chunk sequentially from MinIO
5. The client receives a continuous file stream

---

## Implementation Notes

### Streaming

The service streams data directly to the response using `io.Copy`. Files are not fully loaded into memory, which ensures consistent memory usage regardless of file size.

### Ordering

Chunk ordering is enforced at the database level using `ORDER BY chunk_index`. Object storage is not used to determine order.

### Separation of Concerns

* Handlers manage HTTP concerns only
* Database layer is responsible for metadata
* Object storage layer is responsible for retrieving binary data

---

## Limitations

* No upload endpoint
* No authentication or authorization
* No support for partial content (range requests)
* No data integrity verification (checksum)

---

## Future Work

* Implement upload with chunking
* Add support for HTTP range requests
* Introduce checksum validation
* Add background cleanup jobs
* Improve observability (metrics, structured logging)

---

## License

MIT
