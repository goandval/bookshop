package service

import (
	"context"

	"github.com/yourorg/bookshop/internal/domain"
)

type BookService interface {
	GetByID(ctx context.Context, id int) (*domain.Book, error)
	List(ctx context.Context, categoryIDs []int, limit, offset int) ([]*domain.Book, error)
	Create(ctx context.Context, book *domain.Book) error
	Update(ctx context.Context, book *domain.Book) error
	Delete(ctx context.Context, id int) error
}

type CategoryService interface {
	GetByID(ctx context.Context, id int) (*domain.Category, error)
	List(ctx context.Context) ([]*domain.Category, error)
	Create(ctx context.Context, category *domain.Category) error
	Update(ctx context.Context, category *domain.Category) error
	Delete(ctx context.Context, id int) error
}

type CartService interface {
	GetByUserID(ctx context.Context, userID string) (*domain.Cart, error)
	AddItem(ctx context.Context, userID string, bookID int) error
	RemoveItem(ctx context.Context, userID string, bookID int) error
	Clear(ctx context.Context, userID string) error
	ListItems(ctx context.Context, userID string) ([]*domain.CartItem, error)
}

type OrderService interface {
	Create(ctx context.Context, userID string) (*domain.Order, error)
	ListByUser(ctx context.Context, userID string) ([]*domain.Order, error)
}

type UserService interface {
	GetOrCreate(ctx context.Context, id, email string, isAdmin bool) (*domain.User, error)
}
