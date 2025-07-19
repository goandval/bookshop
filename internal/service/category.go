package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/yourorg/bookshop/internal/domain"
	"github.com/yourorg/bookshop/internal/repository"
)

type CategoryServiceImpl struct {
	categoryRepo repository.CategoryRepository
	bookRepo     repository.BookRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository, bookRepo repository.BookRepository) *CategoryServiceImpl {
	return &CategoryServiceImpl{
		categoryRepo: categoryRepo,
		bookRepo:     bookRepo,
	}
}

func (s *CategoryServiceImpl) GetByID(ctx context.Context, id int) (*domain.Category, error) {
	return s.categoryRepo.GetByID(ctx, id)
}

func (s *CategoryServiceImpl) List(ctx context.Context) ([]*domain.Category, error) {
	return s.categoryRepo.List(ctx)
}

func (s *CategoryServiceImpl) Create(ctx context.Context, category *domain.Category) error {
	if category.Name == "" {
		return fmt.Errorf("name required: %w", errors.New("name required"))
	}
	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return fmt.Errorf("create category: %w", err)
	}
	return nil
}

func (s *CategoryServiceImpl) Update(ctx context.Context, category *domain.Category) error {
	if category.Name == "" {
		return fmt.Errorf("name required: %w", errors.New("name required"))
	}
	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return fmt.Errorf("update category: %w", err)
	}
	return nil
}

func (s *CategoryServiceImpl) Delete(ctx context.Context, id int) error {
	// Найти "без категории"
	noCat, err := s.categoryRepo.GetByName(ctx, "Без категории")
	if err != nil {
		return fmt.Errorf("find 'Без категории': %w", err)
	}
	// Перевести книги в "без категории"
	books, err := s.bookRepo.List(ctx, []int{id}, 10000, 0)
	if err != nil {
		return fmt.Errorf("list books: %w", err)
	}
	for _, b := range books {
		b.CategoryID = noCat.ID
		if err := s.bookRepo.Update(ctx, b); err != nil {
			return fmt.Errorf("move book to 'Без категории': %w", err)
		}
	}
	if err := s.categoryRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete category: %w", err)
	}
	return nil
}
