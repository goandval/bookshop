package domain

import (
	"time"
)

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Book struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Author     string    `json:"author"`
	Year       int       `json:"year"`
	Price      float64   `json:"price"`
	CategoryID int       `json:"category_id"`
	Category   *Category `json:"category,omitempty"`
	Inventory  int       `json:"inventory"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type User struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
}

type Cart struct {
	ID        int        `json:"id"`
	UserID    string     `json:"user_id"`
	Items     []CartItem `json:"items"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type CartItem struct {
	ID         int       `json:"id"`
	CartID     int       `json:"cart_id"`
	BookID     int       `json:"book_id"`
	Book       *Book     `json:"book,omitempty"`
	Quantity   int       `json:"quantity"`
	ReservedAt time.Time `json:"reserved_at"`
}

type Order struct {
	ID        int         `json:"id"`
	UserID    string      `json:"user_id"`
	Items     []OrderItem `json:"items"`
	CreatedAt time.Time   `json:"created_at"`
}

type OrderItem struct {
	ID       int     `json:"id"`
	OrderID  int     `json:"order_id"`
	BookID   int     `json:"book_id"`
	Book     *Book   `json:"book,omitempty"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}
