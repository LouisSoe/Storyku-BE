# Project Name
STORYKU-BE — Story Management Backend

## Introduction
REST API untuk manajemen cerita yang dibangun menggunakan arsitektur bersih (Clean Architecture) dengan bahasa pemrograman Go (Golang) dan database PostgreSQL.

## Table of Contents
- [Introduction](#introduction)
- [Features](#features)
- [Libraries](#libraries)
- [Project Structure](#project-structure)
- [Setup Instructions](#setup-instructions)
- [API / Website URL](#api--website-url)

## Features
- Dashboard
- Story List (dengan pencarian, filter kategori, status, dan pagination)
- Add Story (dengan upload cover image)
- Story Detail
- Edit Story
- Delete Story
- Add Chapter to Story
- Edit Chapter
- Delete Chapter
- Category Management (List, Add, Detail, Edit, Delete)
- Tag Management (List, Add, Detail, Edit, Delete)

## Libraries
- **Echo v4** (`github.com/labstack/echo/v4`) - Web framework ringan dan cepat untuk Go.
- **Godotenv** (`github.com/joho/godotenv`) - Memuat variabel environment dari file `.env`.
- **PostgreSQL Driver** (`github.com/lib/pq`) - Driver database PostgreSQL murni untuk database/sql.
- **Logrus** (`github.com/sirupsen/logrus`) - Logger tersetruktur dengan fitur custom hook (digunakan untuk daily error log).
- **UUID** (`github.com/google/uuid`) - Generator dan parser UUID standar untuk Golang.
- **Go Crypto** (`golang.org/x/crypto`) - Kumpulan package kriptografi tambahan dari Golang.

## Project Structure
- `cmd/` - Titik masuk (entry point) aplikasi (`main.go`).
- `config/` - Pengaturan konfigurasi aplikasi, database, dan logging.
- `core/` - Komponen inti aplikasi berdasarkan Clean Architecture (Domain, Usecase, Repository interfaces).
- `interfaces/` - Implementasi dari layer terluar (Database, HTTP handlers).
- `logs/` - File output log terekam aplikasi.
- `migrations/` - Kumpulan skrip SQL untuk migrasi skema database.
- `pkg/` - Utility package murni yang bisa digunakan kembali (Pagination, Response builders, dll).
- `postman/` - Berisi koleksi Postman untuk pengujian API.
- `routes/` - Tempat pendaftaran rute-rute endpoint API.

## Setup Instructions

### 1. Install Golang
Jika Anda belum menginstall Golang, ikuti langkah berikut:
1. Buka halaman resmi download Golang: [https://go.dev/dl/](https://go.dev/dl/)
2. Unduh *installer* yang sesuai dengan OS yang Anda gunakan (Windows/macOS/Linux).
3. Jalankan *installer* dan ikuti instruksi yang muncul di layar hingga selesai.
4. Buka terminal (atau Command Prompt / PowerShell) baru, lalu ketik `go version`. Jika muncul versi Golang, berarti instalasi berhasil.

### 2. Install PostgreSQL
Jika Anda belum menginstall PostgreSQL:
1. Buka halaman unduhan PostgreSQL: [https://www.postgresql.org/download/](https://www.postgresql.org/download/)
2. Silakan pilih OS Anda dan unduh *installer* (untuk Windows disarankan menggunakan EDB installer).
3. Lakukan instalasi, ingat **password** untuk user `postgres` yang Anda buat saat instalasi (password bawaan biasanya "postgres" atau Anda bisa tentukan sendiri).
4. Pastikan `psql` (PostgreSQL Command Line Tool) sudah masuk ke *Environment Variables / PATH* di OS Anda agar bisa dijalankan langsung dari terminal.

### 3. Setup Project dan Database
1. Buka folder `server` di terminal editor Anda:
   ```bash
   cd server
   ```
2. Salin template environment:
   ```bash
   cp .env.example .env
   ```
   *(Atur konfigurasi password dan username DB di `.env` sesuai dengan instalasi Anda)*
3. Unduh seluruh dependensi Go yang dibutuhkan:
   ```bash
   go mod tidy
   ```
4. Buat database `storyku_db` melalui `psql`:
   ```bash
   psql -U postgres -c "CREATE DATABASE storyku_db;"
   ```
5. Jalankan file migrasi untuk membuat tabel-tabel ke dalam database:
   ```bash
   psql -U postgres -d storyku_db -f migrations/001_create_tables.sql
   ```
6. Jalankan Server:
   - **Mode Manual Runner:**
     ```bash
     go run cmd/main.go
     ```
   - **Mode Build (Production):**
     ```bash
     go build -o build/main.exe ./cmd/main.go
     ./build/main.exe
     ```
   Server akan berjalan secara lokal di port `8080`.

### 4. Menggunakan Postman Collection
Untuk memudahkan pengetesan API backend, ikuti langkah ini:
1. Buka aplikasi **Postman**.
2. Klik tombol **Import** yang berada di bagian kiri atas workspace.
3. Pilih opsi **Upload Files** dan cari serta pilih file `storyku.postman_collection.json` di dalam subdirektori `server/postman/`.
4. Koleksi *endpoints* (request examples) kini sudah masuk dan siap digunakan. Pastikan server lokal Anda sudah berjalan (langkah ke-3).

## API / Website URL
**Backend Base URL:**
[http://103.174.114.118:8081/api/v1](http://103.174.114.118:8081/api/v1)