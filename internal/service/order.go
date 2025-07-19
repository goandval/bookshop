package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/yourorg/bookshop/internal/domain"
	"github.com/yourorg/bookshop/internal/integration"
	"github.com/yourorg/bookshop/internal/repository"
)

type OrderServiceImpl struct {
	orderRepo repository.OrderRepository
	cartRepo  repository.CartRepository
	bookRepo  repository.BookRepository
	kafka     integration.KafkaProducer
	redis     integration.RedisCache
}

func NewOrderService(orderRepo repository.OrderRepository, cartRepo repository.CartRepository, bookRepo repository.BookRepository, kafka integration.KafkaProducer, redis integration.RedisCache) *OrderServiceImpl {
	return &OrderServiceImpl{
		orderRepo: orderRepo,
		cartRepo:  cartRepo,
		bookRepo:  bookRepo,
		kafka:     kafka,
		redis:     redis,
	}
}

func (s *OrderServiceImpl) Create(ctx context.Context, userID string) (*domain.Order, error) {
	items, err := s.ListItems(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get cart items: %w", err)
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("cart is empty: %w", errors.New("cart is empty"))
	}
	var orderItems []domain.OrderItem
	var books []integration.OrderPlacedBook
	for _, item := range items {
		book, err := s.bookRepo.GetByID(ctx, item.BookID)
		if err != nil {
			return nil, fmt.Errorf("book not found: %w", err)
		}
		if book.Inventory < item.Quantity {
			return nil, fmt.Errorf("book out of stock: %d", book.ID)
		}
		orderItems = append(orderItems, domain.OrderItem{BookID: book.ID, Price: book.Price, Quantity: item.Quantity})
		books = append(books, integration.OrderPlacedBook{BookID: book.ID, Quantity: item.Quantity})
	}
	order := &domain.Order{UserID: userID, Items: orderItems}
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("create order: %w", err)
	}
	if err := s.cartRepo.Clear(ctx, userID); err != nil {
		return nil, fmt.Errorf("clear cart: %w", err)
	}
	// Удаляем все резервы пользователя
	for _, item := range items {
		reserveKey := fmt.Sprintf("reserve:%s:%d", userID, item.BookID)
		s.redis.Del(reserveKey)
	}
	if err := s.kafka.PublishOrderPlaced(ctx, order.ID, userID, books); err != nil {
		return nil, fmt.Errorf("publish kafka: %w", err)
	}
	return order, nil
}

func (s *OrderServiceImpl) ListByUser(ctx context.Context, userID string) ([]*domain.Order, error) {
	orders, err := s.orderRepo.ListByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list orders: %w", err)
	}
	return orders, nil
}

func (s *OrderServiceImpl) ListItems(ctx context.Context, userID string) ([]*domain.CartItem, error) {
	items, err := s.cartRepo.ListItems(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list items: %w", err)
	}
	return items, nil
}
