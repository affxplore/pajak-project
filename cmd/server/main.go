package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pajak-backend/config"
	"pajak-backend/database"
	"pajak-backend/internal/handler"
	"pajak-backend/internal/repository"
	"pajak-backend/internal/service"
)

func main() {
	log.Println("==================================================")
	log.Println(" Memulai Suara Rakyat: Pajak Mencekik Backend API ")
	log.Println("==================================================")

	// 1. Inisialisasi Konfigurasi
	cfg := config.LoadConfig()

	// 2. Inisialisasi Koneksi Database PostgreSQL
	db, err := database.ConnectPostgres(cfg)
	if err != nil {
		log.Fatalf("[FATAL] Gagal menghubungkan ke database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("[ERROR] Gagal menutup koneksi database: %v\n", err)
		} else {
			log.Println("[INFO] Koneksi database berhasil ditutup.")
		}
	}()

	// 3. Dependency Injection (Repository -> Service -> Handler)
	commentRepo := repository.NewCommentRepository(db)
	commentServ := service.NewCommentService(commentRepo)
	commentHand := handler.NewCommentHandler(commentServ)

	// 4. Konfigurasi HTTP Router (Go 1.22+ Standard ServeMux)
	mux := http.NewServeMux()

	// Endpoint Healthcheck
	mux.HandleFunc("GET /health", commentHand.HealthCheck)
	mux.HandleFunc("GET /api/health", commentHand.HealthCheck)

	// Endpoint Comments
	mux.HandleFunc("GET /api/comments", commentHand.GetAll)
	mux.HandleFunc("POST /api/comments", commentHand.Create)

	// 5. Bungkus Router dengan Middleware (CORS & Logging)
	rootHandler := handler.LoggingMiddleware(handler.CORSMiddleware(mux))

	// 6. Konfigurasi HTTP Server dengan Timeout yang Aman
	serverAddr := fmt.Sprintf(":%s", cfg.AppPort)
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      rootHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 7. Menjalankan Server secara Asynchronous dalam Goroutine
	go func() {
		log.Printf("[READY] Server berjalan di http://0.0.0.0:%s\n", cfg.AppPort)
		log.Printf("[INFO] API Comments siap menerima request di /api/comments\n")

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[FATAL] Server gagal berjalan: %v\n", err)
		}
	}()

	// 8. Graceful Shutdown (Menunggu sinyal penghentian seperti CTRL+C / SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("\n[SHUTDOWN] Menerima sinyal berhenti, mematikan server secara bertahap...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("[FATAL] Server dipaksa berhenti: %v\n", err)
	}

	log.Println("[SUCCESS] Server Golang berhasil dimatikan dengan aman.")
}

