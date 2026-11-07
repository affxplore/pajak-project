package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"pajak-backend/internal/model"
	"pajak-backend/internal/service"
)

// CommentHandler mengelola HTTP request untuk data kritik/komentar
type CommentHandler struct {
	service service.CommentService
}

// NewCommentHandler membuat instance CommentHandler
func NewCommentHandler(service service.CommentService) *CommentHandler {
	return &CommentHandler{service: service}
}

// Create menangani endpoint POST /api/comments
func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Batasi ukuran request payload maksimal 1 MB untuk keamanan
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	var req model.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Format JSON tidak valid atau ukuran data melebihi batas")
		return
	}

	comment, err := h.service.CreateComment(r.Context(), req)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, http.StatusCreated, model.APIResponse{
		Success: true,
		Message: "Aspirasi/kritik berhasil dipublikasikan",
		Data:    comment,
	})
}

// GetAll menangani endpoint GET /api/comments
func (h *CommentHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Membaca query parameter ?limit=50 (opsional)
	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	comments, total, err := h.service.GetAllComments(r.Context(), limit)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Gagal memuat data kritik: "+err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, model.APIResponse{
		Success: true,
		Message: "Berhasil mengambil daftar kritik",
		Data:    comments,
		Total:   &total,
	})
}

// HealthCheck menangani endpoint GET /health atau GET /api/health
func (h *CommentHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	h.writeJSON(w, http.StatusOK, map[string]string{
		"status":  "healthy",
		"service": "Pajak Mencekik Backend API",
	})
}

// writeJSON adalah fungsi utilitas untuk mengirimkan respons format JSON
func (h *CommentHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, `{"error":"Gagal serialize response JSON"}`, http.StatusInternalServerError)
	}
}

// writeError adalah fungsi utilitas untuk mengirimkan respons error JSON
func (h *CommentHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, model.APIResponse{
		Success: false,
		Message: message,
	})
}

