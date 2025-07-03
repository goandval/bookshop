package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/yourorg/bookshop/internal/domain"
	"github.com/yourorg/bookshop/internal/repository"
)

type BookServiceImpl struct {
	bookRepo     repository.BookRepository
	categoryRepo repository.CategoryRepository
}

func NewBookService(bookRepo repository.BookRepository, categoryRepo repository.CategoryRepository) *BookServiceImpl {
	return &BookServiceImpl{
		bookRepo:     bookRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *BookServiceImpl) GetByID(ctx context.Context, id int) (*domain.Book, error) {
	return s.bookRepo.GetByID(ctx, id)
}

func (s *BookServiceImpl) List(ctx context.Context, categoryIDs []int, limit, offset int) ([]*domain.Book, error) {
	return s.bookRepo.List(ctx, categoryIDs, limit, offset)
}

func (s *BookServiceImpl) Create(ctx context.Context, book *domain.Book) error {
	if book.Inventory < 0 {
		return errors.New("inventory must be >= 0")
	}
	if book.CategoryID == 0 {
		return errors.New("category required")
	}
	return s.bookRepo.Create(ctx, book)
}

func (s *BookServiceImpl) Update(ctx context.Context, book *domain.Book) error {
	// inventory нельзя менять напрямую
	old, err := s.bookRepo.GetByID(ctx, book.ID)
	if err != nil {
		return fmt.Errorf("book not found: %w", err)
	}
	book.Inventory = old.Inventory
	return s.bookRepo.Update(ctx, book)
}

func (s *BookServiceImpl) Delete(ctx context.Context, id int) error {
	return s.bookRepo.Delete(ctx, id)
}
