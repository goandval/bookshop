package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourorg/bookshop/internal/domain"
)

type CartPostgres struct {
	db *pgxpool.Pool
}

func NewCartPostgres(db *pgxpool.Pool) *CartPostgres {
	return &CartPostgres{db: db}
}

func (r *CartPostgres) GetByUserID(ctx context.Context, userID string) (*domain.Cart, error) {
	row := r.db.QueryRow(ctx, `SELECT id, user_id, created_at, updated_at FROM carts WHERE user_id=$1`, userID)
	var c domain.Cart
	if err := row.Scan(&c.ID, &c.UserID, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return nil, fmt.Errorf("get by user: %w", err)
	}
	return &c, nil
}

func (r *CartPostgres) AddItem(ctx context.Context, userID string, bookID int) error {
	cart, err := r.GetByUserID(ctx, userID)
	if err != nil {
		// если корзины нет — создаём
		row := r.db.QueryRow(ctx, `INSERT INTO carts (user_id) VALUES ($1) RETURNING id, created_at, updated_at`, userID)
		cart = &domain.Cart{UserID: userID}
		if err := row.Scan(&cart.ID, &cart.CreatedAt, &cart.UpdatedAt); err != nil {
			return fmt.Errorf("add item: %w", err)
		}
	}
	// Пытаемся увеличить quantity, если книга уже есть
	res, err := r.db.Exec(ctx, `UPDATE cart_items SET quantity = quantity + 1 WHERE cart_id=$1 AND book_id=$2`, cart.ID, bookID)
	n := res.RowsAffected()
	if n == 0 {
		// если не было — вставляем новую строку
		_, err = r.db.Exec(ctx, `INSERT INTO cart_items (cart_id, book_id, quantity) VALUES ($1, $2, 1)`, cart.ID, bookID)
	}
	if err != nil {
		return fmt.Errorf("add item: %w", err)
	}
	return nil
}

func (r *CartPostgres) RemoveItem(ctx context.Context, userID string, bookID int) error {
	cart, err := r.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("remove item: %w", err)
	}
	// Получаем текущий quantity
	var quantity int
	err = r.db.QueryRow(ctx, `SELECT quantity FROM cart_items WHERE cart_id=$1 AND book_id=$2`, cart.ID, bookID).Scan(&quantity)
	if err != nil {
		return fmt.Errorf("get quantity: %w", err)
	}
	if quantity > 1 {
		_, err = r.db.Exec(ctx, `UPDATE cart_items SET quantity = quantity - 1 WHERE cart_id=$1 AND book_id=$2`, cart.ID, bookID)
		if err != nil {
			return fmt.Errorf("remove item: %w", err)
		}
		return nil
	}
	_, err = r.db.Exec(ctx, `DELETE FROM cart_items WHERE cart_id=$1 AND book_id=$2`, cart.ID, bookID)
	if err != nil {
		return fmt.Errorf("remove item: %w", err)
	}
	return nil
}

func (r *CartPostgres) Clear(ctx context.Context, userID string) error {
	cart, err := r.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("clear: %w", err)
	}
	_, err = r.db.Exec(ctx, `DELETE FROM cart_items WHERE cart_id=$1`, cart.ID)
	if err != nil {
		return fmt.Errorf("clear: %w", err)
	}
	return nil
}

func (r *CartPostgres) ListItems(ctx context.Context, userID string) ([]*domain.CartItem, error) {
	cart, err := r.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list items: %w", err)
	}
	rows, err := r.db.Query(ctx, `SELECT id, cart_id, book_id, quantity, reserved_at FROM cart_items WHERE cart_id=$1`, cart.ID)
	if err != nil {
		return nil, fmt.Errorf("list items: %w", err)
	}
	defer rows.Close()
	var items []*domain.CartItem
	for rows.Next() {
		var it domain.CartItem
		if err := rows.Scan(&it.ID, &it.CartID, &it.BookID, &it.Quantity, &it.ReservedAt); err != nil {
			return nil, fmt.Errorf("scan item: %w", err)
		}
		items = append(items, &it)
	}
	return items, nil
}

func (r *CartPostgres) GetItemQuantity(ctx context.Context, userID string, bookID int) (int, error) {
	cart, err := r.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	var quantity int
	err = r.db.QueryRow(ctx, `SELECT quantity FROM cart_items WHERE cart_id=$1 AND book_id=$2`, cart.ID, bookID).Scan(&quantity)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("get item quantity: %w", err)
	}
	return quantity, nil
}
