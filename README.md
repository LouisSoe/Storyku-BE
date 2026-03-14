# STORYKU-BE — Story Management Backend

REST API untuk manajemen cerita — Go, PostgreSQL.

## Setup
```bash
go mod tidy

psql -U postgres -c "CREATE DATABASE storyku_db;"
psql -U postgres -d storyku_db -f migrations/001_create_tables.sql

# Development (hot-reload)
air

# Production
go build -o bin/storyku-be ./cmd/main.go && ./bin/storyku-be
```

## Endpoints

| Method | URL | Deskripsi |
|--------|-----|-----------|
| GET | `/api/v1/stories` | List + search + filter |
| POST | `/api/v1/stories` | Tambah story |
| GET | `/api/v1/stories/:id` | Detail story |
| PUT | `/api/v1/stories/:id` | Edit story |
| DELETE | `/api/v1/stories/:id` | Hapus story |
| POST | `/api/v1/stories/:id/chapters` | Tambah chapter |
| PUT | `/api/v1/stories/:id/chapters/:cid` | Edit chapter |
| DELETE | `/api/v1/stories/:id/chapters/:cid` | Hapus chapter |

## Query Params (GET /api/v1/stories)

`search`, `category` (Financial/Technology/Health), `status` (publish/draft), `page`, `limit`

## Upload Cover

`multipart/form-data`, field `cover`. Format: jpg/jpeg/png/webp, maks 5MB.

## Tests
```bash
go test ./core/usecase/... -v
```