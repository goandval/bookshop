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

	userID := "user-1"
	// cartItems := []*domain.CartItem{{ID: 1, BookID: 42}}

	bookRepo.On("GetByID", mock.Anything, 42).Return(&domain.Book{ID: 42, Inventory: 1, Price: 10.0}, nil)
	orderRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Order")).Return(nil)
	cartRepo.On("Clear", mock.Anything, userID).Return(nil)
	cartRepo.On("ListItems", mock.Anything, userID).Return([]*domain.CartItem{{ID: 1, BookID: 42, Quantity: 1}}, nil)
	kafka.On("PublishOrderPlaced", mock.Anything, mock.Anything, userID, []integration.OrderPlacedBook{{BookID: 42, Quantity: 1}}).Return(nil)

	svc := &OrderServiceImpl{orderRepo, cartRepo, bookRepo, kafka}
	res, err := svc.Create(context.Background(), userID)
	require.NoError(t, err)
	assert.NotNil(t, res)
	orderRepo.AssertExpectations(t)
	cartRepo.AssertExpectations(t)
	kafka.AssertExpectations(t)
}
