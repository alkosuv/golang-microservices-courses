package service

import (
	"context"
	"github.com/alkosuv/golang-microservices-courses/docker/storage/internal/model"
	"log/slog"
)

type repository interface {
	GetBooks(ctx context.Context) ([]*model.Book, error)
	AddBook(ctx context.Context, book *model.Book) error
}

type Service struct {
	repository repository
}

func New(repository repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetBooks(ctx context.Context) ([]*model.Book, error) {
	books, err := s.repository.GetBooks(ctx)
	if err != nil {
		slog.Error("failed to get books: %v", err)
		return nil, err
	}

	return books, nil
}

func (s *Service) AddBook(ctx context.Context, book *model.Book) error {
	if err := s.repository.AddBook(ctx, book); err != nil {
		slog.Error("error adding book: %v", err)
		return err
	}

	return nil
}
