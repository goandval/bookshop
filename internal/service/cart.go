package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/yourorg/bookshop/internal/domain"
	"github.com/yourorg/bookshop/internal/integration"
	"github.com/yourorg/bookshop/internal/repository"
)

type CartServiceImpl struct {
	cartRepo repository.CartRepository
	bookRepo repository.BookRepository
	redis    integration.RedisCache
}

func NewCartService(cartRepo repository.CartRepository, bookRepo repository.BookRepository, redis integration.RedisCache) *CartServiceImpl {
	return &CartServiceImpl{
		cartRepo: cartRepo,
		bookRepo: bookRepo,
		redis:    redis,
	}
}

func (s *CartServiceImpl) GetByUserID(ctx context.Context, userID string) (*domain.Cart, error) {
	// TODO: реализовать
	return nil, nil
}

func (s *CartServiceImpl) AddItem(ctx context.Context, userID string, bookID int) error {
	book, err := s.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		return fmt.Errorf("book not found: %w", err)
	}
	if book.Inventory <= 0 {
		return errors.New("book is out of stock")
	}
	// Проверка резервирования через Redis
	redisKey := "reserve:book:" + strconv.Itoa(bookID)
	val, err := s.redis.Get(redisKey)
	if err == nil && val != "" {
		return errors.New("book is already reserved")
	}
	// Добавление в корзину
	if err := s.cartRepo.AddItem(ctx, userID, bookID); err != nil {
		return fmt.Errorf("add to cart: %w", err)
	}
	// Резервируем в Redis на 30 минут
	if err := s.redis.Set(redisKey, userID, 1800); err != nil {
		return fmt.Errorf("reserve in redis: %w", err)
	}
	return nil
}

func (s *CartServiceImpl) RemoveItem(ctx context.Context, userID string, bookID int) error {
	// TODO: реализовать
	return nil
}

func (s *CartServiceImpl) Clear(ctx context.Context, userID string) error {
	// TODO: реализовать
	return nil
}

func (s *CartServiceImpl) ListItems(ctx context.Context, userID string) ([]*domain.CartItem, error) {
	// TODO: реализовать
	return nil, nil
}
