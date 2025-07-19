package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/yourorg/bookshop/internal/domain"
	"github.com/yourorg/bookshop/internal/integration"
	"github.com/yourorg/bookshop/internal/repository"
	"golang.org/x/exp/slog"
)

type CartServiceImpl struct {
	cartRepo repository.CartRepository
	bookRepo repository.BookRepository
	redis    integration.RedisCache
	Logger   *slog.Logger
}

func NewCartService(cartRepo repository.CartRepository, bookRepo repository.BookRepository, redis integration.RedisCache, logger *slog.Logger) *CartServiceImpl {
	return &CartServiceImpl{
		cartRepo: cartRepo,
		bookRepo: bookRepo,
		redis:    redis,
		Logger:   logger,
	}
}

func (s *CartServiceImpl) GetByUserID(ctx context.Context, userID string) (*domain.Cart, error) {
	cart, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	items, err := s.cartRepo.ListItems(ctx, userID)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		book, err := s.bookRepo.GetByID(ctx, item.BookID)
		if err == nil {
			item.Book = book
		}
	}
	cart.Items = nil
	for _, item := range items {
		if item != nil {
			cart.Items = append(cart.Items, *item)
		}
	}
	return cart, nil
}

func (s *CartServiceImpl) AddItem(ctx context.Context, userID string, bookID int) error {
	book, err := s.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		return fmt.Errorf("book not found: %w", err)
	}
	if book.Inventory <= 0 {
		return fmt.Errorf("book is out of stock: %w", errors.New("book is out of stock"))
	}
	quantity, err := s.cartRepo.GetItemQuantity(ctx, userID, bookID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			quantity = 0
		} else {
			return fmt.Errorf("get item quantity: %w", err)
		}
	}
	if quantity >= book.Inventory {
		return fmt.Errorf("not enough books in stock: %w", errors.New("not enough books in stock"))
	}
	// inventory НЕ уменьшаем! Только добавляем в корзину
	if err := s.cartRepo.AddItem(ctx, userID, bookID); err != nil {
		return fmt.Errorf("add to cart: %w", err)
	}
	return nil
}

func (s *CartServiceImpl) RemoveItem(ctx context.Context, userID string, bookID int) error {
	if err := s.cartRepo.RemoveItem(ctx, userID, bookID); err != nil {
		return fmt.Errorf("remove item: %w", err)
	}
	return nil
}

func (s *CartServiceImpl) Clear(ctx context.Context, userID string) error {
	if err := s.cartRepo.Clear(ctx, userID); err != nil {
		return fmt.Errorf("clear cart: %w", err)
	}
	return nil
}

func (s *CartServiceImpl) ListItems(ctx context.Context, userID string) ([]*domain.CartItem, error) {
	items, err := s.cartRepo.ListItems(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list items: %w", err)
	}
	for _, item := range items {
		book, err := s.bookRepo.GetByID(ctx, item.BookID)
		if err == nil {
			item.Book = book
		}
	}
	return items, nil
}
