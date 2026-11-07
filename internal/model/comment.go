package model

import (
	"errors"
	"strings"
	"time"
)

// Comment merepresentasikan entitas kritik/komentar dalam database
type Comment struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateCommentRequest menampung payload JSON yang dikirimkan oleh frontend
type CreateCommentRequest struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	Comment  string `json:"comment"`
}

// Validate memvalidasi data masukan dari request form frontend
func (r *CreateCommentRequest) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	r.Category = strings.TrimSpace(r.Category)
	r.Comment = strings.TrimSpace(r.Comment)

	if r.Name == "" {
		return errors.New("nama atau inisial tidak boleh kosong")
	}
	if len(r.Name) > 100 {
		return errors.New("nama maksimal 100 karakter")
	}

	if r.Category == "" {
		r.Category = "Umum"
	}
	if len(r.Category) > 50 {
		return errors.New("kategori maksimal 50 karakter")
	}

	if r.Comment == "" {
		return errors.New("isi komentar/aspirasi tidak boleh kosong")
	}
	if len(r.Comment) > 2000 {
		return errors.New("komentar maksimal 2000 karakter")
	}

	return nil
}

// APIResponse adalah format standar respons JSON untuk klien
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Total   *int64      `json:"total,omitempty"`
}

