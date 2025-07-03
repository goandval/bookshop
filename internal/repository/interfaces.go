package repository

import (
	"context"

	"github.com/yourorg/bookshop/internal/domain"
)

type BookRepository interface {
	GetByID(ctx context.Context, id int) (*domain.Book, error)
	List(ctx context.Context, categoryIDs []int, limit, offset int) ([]*domain.Book, error)
	Create(ctx context.Context, book *domain.Book) error
	Update(ctx context.Context, book *domain.Book) error
	Delete(ctx context.Context, id int) error
}

type CategoryRepository interface {
	GetByID(ctx context.Context, id int) (*domain.Category, error)
	List(ctx context.Context) ([]*domain.Category, error)
	Create(ctx context.Context, category *domain.Category) error
	Update(ctx context.Context, category *domain.Category) error
	Delete(ctx context.Context, id int) error
	GetByName(ctx context.Context, name string) (*domain.Category, error)
}

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	CreateIfNotExists(ctx context.Context, user *domain.User) error
}

type CartRepository interface {
	GetByUserID(ctx context.Context, userID string) (*domain.Cart, error)
	AddItem(ctx context.Context, userID string, bookID int) error
	RemoveItem(ctx context.Context, userID string, bookID int) error
	Clear(ctx context.Context, userID string) error
	ListItems(ctx context.Context, userID string) ([]*domain.CartItem, error)
}

type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	ListByUser(ctx context.Context, userID string) ([]*domain.Order, error)
}
