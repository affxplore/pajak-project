package handler

import (
	"log"
	"net/http"
	"time"
)

// CORSMiddleware menangani Cross-Origin Resource Sharing agar frontend HTML dapat memanggil API
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mengizinkan request dari semua origin (dapat disesuaikan jika ingin membatasi domain tertentu)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Jika request adalah HTTP OPTIONS (Preflight check dari browser), langsung balas dengan 204 No Content
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware mencatat waktu proses setiap HTTP request
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[%s] %s | Durasi: %v\n", r.Method, r.URL.Path, time.Since(start))
	})
}

