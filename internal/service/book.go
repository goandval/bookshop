package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/yourorg/bookshop/internal/domain"
	"github.com/yourorg/bookshop/internal/integration"
	"github.com/yourorg/bookshop/internal/repository"
)

type BookServiceImpl struct {
	bookRepo     repository.BookRepository
	categoryRepo repository.CategoryRepository
	redis        integration.RedisCache
}

func NewBookService(bookRepo repository.BookRepository, categoryRepo repository.CategoryRepository, redis integration.RedisCache) *BookServiceImpl {
	return &BookServiceImpl{
		bookRepo:     bookRepo,
		categoryRepo: categoryRepo,
		redis:        redis,
	}
}

func (s *BookServiceImpl) GetByID(ctx context.Context, id int) (*domain.Book, error) {
	return s.bookRepo.GetByID(ctx, id)
}

func (s *BookServiceImpl) List(ctx context.Context, categoryIDs []int, limit, offset int) ([]*domain.Book, error) {
	if limit != 100 || offset != 0 {
		return s.bookRepo.List(ctx, categoryIDs, limit, offset)
	}
	key := "books:all"
	if len(categoryIDs) == 1 {
		key = "books:cat:" + fmt.Sprint(categoryIDs[0])
	}
	if cached, err := s.redis.Get(key); err == nil && cached != "" {
		var books []*domain.Book
		if err := json.Unmarshal([]byte(cached), &books); err == nil {
			return books, nil
		}
	}
	books, err := s.bookRepo.List(ctx, categoryIDs, limit, offset)
	if err == nil {
		if data, err := json.Marshal(books); err == nil {
			s.redis.Set(key, string(data), 300) // 5 минут
		}
	}
	return books, err
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
	s.redis.Del("books:all")
	s.redis.Del("books:cat:" + fmt.Sprint(book.CategoryID))
	return nil
}

func (s *BookServiceImpl) Update(ctx context.Context, book *domain.Book) error {
	old, err := s.bookRepo.GetByID(ctx, book.ID)
	if err != nil {
		return fmt.Errorf("book not found: %w", err)
	}
	book.Inventory = old.Inventory
	if err := s.bookRepo.Update(ctx, book); err != nil {
		return fmt.Errorf("update book: %w", err)
	}
	s.redis.Del("books:all")
	s.redis.Del("books:cat:" + fmt.Sprint(book.CategoryID))
	return nil
}

func (s *BookServiceImpl) Delete(ctx context.Context, id int) error {
	book, _ := s.bookRepo.GetByID(ctx, id)
	if err := s.bookRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete book: %w", err)
	}
	s.redis.Del("books:all")
	if book != nil {
		s.redis.Del("books:cat:" + fmt.Sprint(book.CategoryID))
	}
	return nil
}
