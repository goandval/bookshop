package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourorg/bookshop/internal/domain"
)

type OrderPostgres struct {
	db *pgxpool.Pool
}

func NewOrderPostgres(db *pgxpool.Pool) *OrderPostgres {
	return &OrderPostgres{db: db}
}

func (r *OrderPostgres) Create(ctx context.Context, order *domain.Order) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)
	row := tx.QueryRow(ctx, `INSERT INTO orders (user_id) VALUES ($1) RETURNING id, created_at`, order.UserID)
	if err := row.Scan(&order.ID, &order.CreatedAt); err != nil {
		return fmt.Errorf("insert order: %w", err)
	}
	for _, item := range order.Items {
		row := tx.QueryRow(ctx, `INSERT INTO order_items (order_id, book_id, price, quantity) VALUES ($1, $2, $3, $4) RETURNING id`, order.ID, item.BookID, item.Price, item.Quantity)
		if err := row.Scan(&item.ID); err != nil {
			return fmt.Errorf("insert item: %w", err)
		}
		// Списываем остаток
		_, err := tx.Exec(ctx, `UPDATE books SET inventory = inventory - 1 WHERE id=$1 AND inventory > 0`, item.BookID)
		if err != nil {
			return fmt.Errorf("update inventory: %w", err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	return nil
}

func (r *OrderPostgres) ListByUser(ctx context.Context, userID string) ([]*domain.Order, error) {
	rows, err := r.db.Query(ctx, `SELECT id, created_at FROM orders WHERE user_id=$1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("list by user: %w", err)
	}
	defer rows.Close()
	var orders []*domain.Order
	for rows.Next() {
		var o domain.Order
		o.UserID = userID
		if err := rows.Scan(&o.ID, &o.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan order: %w", err)
		}
		// Получаем order_items
		itemRows, err := r.db.Query(ctx, `SELECT id, order_id, book_id, price, quantity FROM order_items WHERE order_id=$1`, o.ID)
		if err != nil {
			return nil, fmt.Errorf("get order items: %w", err)
		}
		for itemRows.Next() {
			var it domain.OrderItem
			if err := itemRows.Scan(&it.ID, &it.OrderID, &it.BookID, &it.Price, &it.Quantity); err != nil {
				itemRows.Close()
				return nil, fmt.Errorf("scan item: %w", err)
			}
			o.Items = append(o.Items, it)
		}
		itemRows.Close()
		orders = append(orders, &o)
	}
	return orders, nil
}
