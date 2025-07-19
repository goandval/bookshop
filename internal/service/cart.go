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
		return errors.New("book is out of stock")
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
		return errors.New("not enough books in stock")
	}
	// inventory НЕ уменьшаем! Только добавляем в корзину
	if err := s.cartRepo.AddItem(ctx, userID, bookID); err != nil {
		return fmt.Errorf("add to cart: %w", err)
	}
	return nil
}

func (s *CartServiceImpl) RemoveItem(ctx context.Context, userID string, bookID int) error {
	return s.cartRepo.RemoveItem(ctx, userID, bookID)
}

func (s *CartServiceImpl) Clear(ctx context.Context, userID string) error {
	return s.cartRepo.Clear(ctx, userID)
}

func (s *CartServiceImpl) ListItems(ctx context.Context, userID string) ([]*domain.CartItem, error) {
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
	return items, nil
}
