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
}

func NewOrderService(orderRepo repository.OrderRepository, cartRepo repository.CartRepository, bookRepo repository.BookRepository, kafka integration.KafkaProducer) *OrderServiceImpl {
	return &OrderServiceImpl{
		orderRepo: orderRepo,
		cartRepo:  cartRepo,
		bookRepo:  bookRepo,
		kafka:     kafka,
	}
}

func (s *OrderServiceImpl) Create(ctx context.Context, userID string) (*domain.Order, error) {
	items, err := s.ListItems(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get cart items: %w", err)
	}
	if len(items) == 0 {
		return nil, errors.New("cart is empty")
	}
	var orderItems []domain.OrderItem
	var bookIDs []int
	for _, item := range items {
		book, err := s.bookRepo.GetByID(ctx, item.BookID)
		if err != nil {
			return nil, fmt.Errorf("book not found: %w", err)
		}
		if book.Inventory <= 0 {
			return nil, fmt.Errorf("book out of stock: %d", book.ID)
		}
		orderItems = append(orderItems, domain.OrderItem{BookID: book.ID, Price: book.Price})
		bookIDs = append(bookIDs, book.ID)
	}
	// Атомарное списание остатков и создание заказа (эмулируем транзакцию)
	order := &domain.Order{UserID: userID, Items: orderItems}
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("create order: %w", err)
	}
	if err := s.cartRepo.Clear(ctx, userID); err != nil {
		return nil, fmt.Errorf("clear cart: %w", err)
	}
	if err := s.kafka.PublishOrderPlaced(ctx, order.ID, userID, bookIDs); err != nil {
		return nil, fmt.Errorf("publish kafka: %w", err)
	}
	return order, nil
}

func (s *OrderServiceImpl) ListByUser(ctx context.Context, userID string) ([]*domain.Order, error) {
	// TODO: реализовать
	return nil, nil
}

func (s *OrderServiceImpl) ListItems(ctx context.Context, userID string) ([]*domain.CartItem, error) {
	return s.cartRepo.ListItems(ctx, userID)
}
