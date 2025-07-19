package service

import (
	"context"
	"testing"

	"io"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yourorg/bookshop/internal/domain"
	"github.com/yourorg/bookshop/internal/mocks"
	"golang.org/x/exp/slog"
)

func TestCartService_AddItem_Success(t *testing.T) {
	cartRepo := new(mocks.CartRepository)
	bookRepo := new(mocks.BookRepository)
	redis := new(mocks.RedisCache)

	userID := "user-1"
	bookID := 42

	bookRepo.On("GetByID", mock.Anything, bookID).Return(&domain.Book{ID: bookID, Inventory: 1}, nil)
	cartRepo.On("AddItem", mock.Anything, userID, bookID).Return(nil)
	cartRepo.On("GetItemQuantity", mock.Anything, userID, bookID).Return(0, nil)

	svc := &CartServiceImpl{cartRepo, bookRepo, redis, slog.New(slog.NewTextHandler(io.Discard, nil))}
	err := svc.AddItem(context.Background(), userID, bookID)
	require.NoError(t, err)
	bookRepo.AssertExpectations(t)
	cartRepo.AssertExpectations(t)
}

func TestCartService_AddItem_AlreadyReserved(t *testing.T) {
	cartRepo := new(mocks.CartRepository)
	bookRepo := new(mocks.BookRepository)
	redis := new(mocks.RedisCache)

	userID := "user-1"
	bookID := 42

	bookRepo.On("GetByID", mock.Anything, bookID).Return(&domain.Book{ID: bookID, Inventory: 1}, nil)
	cartRepo.On("GetItemQuantity", mock.Anything, userID, bookID).Return(1, nil)

	svc := &CartServiceImpl{cartRepo, bookRepo, redis, slog.New(slog.NewTextHandler(io.Discard, nil))}
	err := svc.AddItem(context.Background(), userID, bookID)
	require.Error(t, err)
	require.Equal(t, "not enough books in stock: not enough books in stock", err.Error())
	cartRepo.AssertExpectations(t)
}
