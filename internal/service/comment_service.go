package service

import (
	"context"
	"fmt"

	"pajak-backend/internal/model"
	"pajak-backend/internal/repository"
)

// CommentService mendefinisikan interface business logic untuk kritik/komentar
type CommentService interface {
	CreateComment(ctx context.Context, req model.CreateCommentRequest) (*model.Comment, error)
	GetAllComments(ctx context.Context, limit int) ([]model.Comment, int64, error)
}

// commentService adalah implementasi dari CommentService
type commentService struct {
	repo repository.CommentRepository
}

// NewCommentService membuat instance baru CommentService
func NewCommentService(repo repository.CommentRepository) CommentService {
	return &commentService{repo: repo}
}

// CreateComment memvalidasi data masukan dan menyimpannya melalui repository
func (s *commentService) CreateComment(ctx context.Context, req model.CreateCommentRequest) (*model.Comment, error) {
	// Validasi business rules
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validasi gagal: %w", err)
	}

	comment := &model.Comment{
		Name:     req.Name,
		Category: req.Category,
		Comment:  req.Comment,
	}

	if err := s.repo.Create(ctx, comment); err != nil {
		return nil, fmt.Errorf("service: gagal menambahkan komentar: %w", err)
	}

	return comment, nil
}

// GetAllComments mengambil daftar kritik beserta total seluruh kritik di database
func (s *commentService) GetAllComments(ctx context.Context, limit int) ([]model.Comment, int64, error) {
	comments, err := s.repo.FindAll(ctx, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("service: gagal mengambil daftar komentar: %w", err)
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("service: gagal mengambil total komentar: %w", err)
	}

	return comments, total, nil
}

