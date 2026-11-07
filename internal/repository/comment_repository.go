package repository

import (
	"context"
	"database/sql"
	"fmt"

	"pajak-backend/internal/model"
)

// CommentRepository mendefinisikan kontrak operasi database untuk entitas Comment
type CommentRepository interface {
	Create(ctx context.Context, comment *model.Comment) error
	FindAll(ctx context.Context, limit int) ([]model.Comment, error)
	Count(ctx context.Context) (int64, error)
}

// commentRepository adalah implementasi CommentRepository menggunakan PostgreSQL (*sql.DB)
type commentRepository struct {
	db *sql.DB
}

// NewCommentRepository membuat instance baru CommentRepository
func NewCommentRepository(db *sql.DB) CommentRepository {
	return &commentRepository{db: db}
}

// Create menyimpan data kritik baru ke dalam tabel comments
func (r *commentRepository) Create(ctx context.Context, comment *model.Comment) error {
	query := `
		INSERT INTO comments (name, category, comment)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		comment.Name,
		comment.Category,
		comment.Comment,
	).Scan(&comment.ID, &comment.CreatedAt)

	if err != nil {
		return fmt.Errorf("repository: gagal menyimpan kritik: %w", err)
	}

	return nil
}

// FindAll mengambil daftar kritik yang diurutkan dari yang terbaru (created_at DESC)
func (r *commentRepository) FindAll(ctx context.Context, limit int) ([]model.Comment, error) {
	if limit <= 0 || limit > 100 {
		limit = 50 // Default limit 50 item
	}

	query := `
		SELECT id, name, category, comment, created_at
		FROM comments
		ORDER BY created_at DESC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("repository: gagal mengambil daftar kritik: %w", err)
	}
	defer rows.Close()

	var comments []model.Comment
	for rows.Next() {
		var c model.Comment
		if err := rows.Scan(&c.ID, &c.Name, &c.Category, &c.Comment, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("repository: gagal scan data kritik: %w", err)
		}
		comments = append(comments, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: error iterasi rows: %w", err)
	}

	// Kembalikan slice kosong (bukan nil) jika belum ada data agar JSON menghasilkan []
	if comments == nil {
		comments = []model.Comment{}
	}

	return comments, nil
}

// Count menghitung total jumlah kritik yang tersimpan di database
func (r *commentRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM comments`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("repository: gagal menghitung total komentar: %w", err)
	}

	return count, nil
}

