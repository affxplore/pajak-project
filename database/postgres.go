package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"pajak-backend/config"

	_ "github.com/lib/pq" // Driver PostgreSQL
)

// ConnectPostgres membuka koneksi ke PostgreSQL dan mengatur connection pool
func ConnectPostgres(cfg *config.Config) (*sql.DB, error) {
	dsn := cfg.GetDSN()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka koneksi ke database: %w", err)
	}

	// Konfigurasi Database Connection Pool untuk performa optimal
	db.SetMaxOpenConns(25)                 // Maksimum koneksi aktif simultan
	db.SetMaxIdleConns(10)                 // Maksimum koneksi idle yang ditahan
	db.SetConnMaxLifetime(15 * time.Minute) // Waktu hidup maksimum koneksi
	db.SetConnMaxIdleTime(5 * time.Minute)  // Batas waktu koneksi idle sebelum ditutup

	// Uji koneksi dengan Ping
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("gagal melakukan ping ke database PostgreSQL: %w", err)
	}

	log.Printf("[SUCCESS] Berhasil terhubung ke database PostgreSQL (%s:%s/%s)\n",
		cfg.DBHost, cfg.DBPort, cfg.DBName)

	return db, nil
}

