package repository

import (
	"context"

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
		return nil, err
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
			return err
		}
	}
	_, err = r.db.Exec(ctx, `INSERT INTO cart_items (cart_id, book_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, cart.ID, bookID)
	return err
}

func (r *CartPostgres) RemoveItem(ctx context.Context, userID string, bookID int) error {
	cart, err := r.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, `DELETE FROM cart_items WHERE cart_id=$1 AND book_id=$2`, cart.ID, bookID)
	return err
}

func (r *CartPostgres) Clear(ctx context.Context, userID string) error {
	cart, err := r.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, `DELETE FROM cart_items WHERE cart_id=$1`, cart.ID)
	return err
}

func (r *CartPostgres) ListItems(ctx context.Context, userID string) ([]*domain.CartItem, error) {
	cart, err := r.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Query(ctx, `SELECT id, cart_id, book_id, reserved_at FROM cart_items WHERE cart_id=$1`, cart.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*domain.CartItem
	for rows.Next() {
		var it domain.CartItem
		if err := rows.Scan(&it.ID, &it.CartID, &it.BookID, &it.ReservedAt); err != nil {
			return nil, err
		}
		items = append(items, &it)
	}
	return items, nil
}
