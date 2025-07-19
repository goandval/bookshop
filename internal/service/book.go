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
		return fmt.Errorf("inventory must be >= 0: %w", errors.New("inventory must be >= 0"))
	}
	if book.CategoryID == 0 {
		return fmt.Errorf("category required: %w", errors.New("category required"))
	}
	if err := s.bookRepo.Create(ctx, book); err != nil {
		return fmt.Errorf("create book: %w", err)
	}
	return nil
}

func (s *BookServiceImpl) Update(ctx context.Context, book *domain.Book) error {
	// inventory нельзя менять напрямую
	old, err := s.bookRepo.GetByID(ctx, book.ID)
	if err != nil {
		return fmt.Errorf("book not found: %w", err)
	}
	book.Inventory = old.Inventory
	if err := s.bookRepo.Update(ctx, book); err != nil {
		return fmt.Errorf("update book: %w", err)
	}
	return nil
}

func (s *BookServiceImpl) Delete(ctx context.Context, id int) error {
	if err := s.bookRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete book: %w", err)
	}
	return nil
}
