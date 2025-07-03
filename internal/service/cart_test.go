package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yourorg/bookshop/internal/domain"
	"github.com/yourorg/bookshop/internal/mocks"
)

func TestCartService_AddItem_Success(t *testing.T) {
	cartRepo := new(mocks.CartRepository)
	bookRepo := new(mocks.BookRepository)
	redis := new(mocks.RedisCache)

	userID := "user-1"
	bookID := 42

	bookRepo.On("GetByID", mock.Anything, bookID).Return(&domain.Book{ID: bookID, Inventory: 1}, nil)
	redis.On("Get", mock.Anything).Return("", errors.New("redis: nil")) // not reserved
	redis.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	cartRepo.On("AddItem", mock.Anything, userID, bookID).Return(nil)

	svc := &CartServiceImpl{cartRepo, bookRepo, redis}
	err := svc.AddItem(context.Background(), userID, bookID)
	assert.NoError(t, err)
	bookRepo.AssertExpectations(t)
	redis.AssertExpectations(t)
	cartRepo.AssertExpectations(t)
}

func TestCartService_AddItem_AlreadyReserved(t *testing.T) {
	cartRepo := new(mocks.CartRepository)
	bookRepo := new(mocks.BookRepository)
	redis := new(mocks.RedisCache)

	userID := "user-1"
	bookID := 42

	bookRepo.On("GetByID", mock.Anything, bookID).Return(&domain.Book{ID: bookID, Inventory: 1}, nil)
	redis.On("Get", mock.Anything).Return("reserved", nil) // already reserved

	svc := &CartServiceImpl{cartRepo, bookRepo, redis}
	err := svc.AddItem(context.Background(), userID, bookID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "reserved")
}
