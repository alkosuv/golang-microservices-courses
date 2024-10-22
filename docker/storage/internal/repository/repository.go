package repository

import (
	"context"
	"github.com/alkosuv/golang-microservices-courses/docker/storage/internal/model"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetBooks(ctx context.Context) ([]*model.Book, error) {
	sql := `SELECT * FROM books`

	var books []*model.Book
	if err := pgxscan.Select(ctx, r.db, &books, sql); err != nil {
		return nil, err
	}

	return books, nil
}

func (r *Repository) AddBook(ctx context.Context, book *model.Book) error {
	sql := `INSERT INTO books (title, author, pages) VALUES ($1, $2, $3)`

	if _, err := r.db.Exec(ctx, sql, book.Title, book.Author, book.Pages); err != nil {
		return err
	}

	return nil
}
