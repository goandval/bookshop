package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yourorg/bookshop/internal/domain"
	"github.com/yourorg/bookshop/internal/integration"
	"github.com/yourorg/bookshop/internal/mocks"
)

func TestOrderService_Create_Success(t *testing.T) {
	orderRepo := new(mocks.OrderRepository)
	cartRepo := new(mocks.CartRepository)
	bookRepo := new(mocks.BookRepository)
	kafka := new(mocks.KafkaProducer)
	redis := new(mocks.RedisCache)

	userID := "user-1"
	// cartItems := []*domain.CartItem{{ID: 1, BookID: 42}}

	bookRepo.On("GetByID", mock.Anything, 42).Return(&domain.Book{ID: 42, Inventory: 1, Price: 10.0}, nil)
	orderRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Order")).Return(nil)
	cartRepo.On("Clear", mock.Anything, userID).Return(nil)
	cartRepo.On("ListItems", mock.Anything, userID).Return([]*domain.CartItem{{ID: 1, BookID: 42, Quantity: 1}}, nil)
	kafka.On("PublishOrderPlaced", mock.Anything, mock.Anything, userID, []integration.OrderPlacedBook{{BookID: 42, Quantity: 1}}).Return(nil)
	redis.On("Del", "reserve:user-1:42").Return(nil)

	svc := &OrderServiceImpl{orderRepo, cartRepo, bookRepo, kafka, redis}
	res, err := svc.Create(context.Background(), userID)
	require.NoError(t, err)
	assert.NotNil(t, res)
	orderRepo.AssertExpectations(t)
	cartRepo.AssertExpectations(t)
	kafka.AssertExpectations(t)
	redis.AssertExpectations(t)
}

func TestOrderService_Create_RemovesReserves(t *testing.T) {
	orderRepo := new(mocks.OrderRepository)
	cartRepo := new(mocks.CartRepository)
	bookRepo := new(mocks.BookRepository)
	kafka := new(mocks.KafkaProducer)
	redis := new(mocks.RedisCache)

	userID := "user-1"
	items := []*domain.CartItem{{BookID: 42, Quantity: 1}, {BookID: 43, Quantity: 2}}

	bookRepo.On("GetByID", mock.Anything, 42).Return(&domain.Book{ID: 42, Inventory: 10, Price: 10.0}, nil)
	bookRepo.On("GetByID", mock.Anything, 43).Return(&domain.Book{ID: 43, Inventory: 10, Price: 20.0}, nil)
	orderRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Order")).Return(nil)
	cartRepo.On("Clear", mock.Anything, userID).Return(nil)
	cartRepo.On("ListItems", mock.Anything, userID).Return(items, nil)
	kafka.On("PublishOrderPlaced", mock.Anything, mock.Anything, userID, mock.Anything).Return(nil)
	redis.On("Del", "reserve:user-1:42").Return(nil)
	redis.On("Del", "reserve:user-1:43").Return(nil)

	svc := &OrderServiceImpl{orderRepo, cartRepo, bookRepo, kafka, redis}
	_, err := svc.Create(context.Background(), userID)
	require.NoError(t, err)
	redis.AssertExpectations(t)
}
