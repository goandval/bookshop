package repository

import (
	"context"

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
		return err
	}
	defer tx.Rollback(ctx)
	row := tx.QueryRow(ctx, `INSERT INTO orders (user_id) VALUES ($1) RETURNING id, created_at`, order.UserID)
	if err := row.Scan(&order.ID, &order.CreatedAt); err != nil {
		return err
	}
	for _, item := range order.Items {
		row := tx.QueryRow(ctx, `INSERT INTO order_items (order_id, book_id, price) VALUES ($1, $2, $3) RETURNING id`, order.ID, item.BookID, item.Price)
		if err := row.Scan(&item.ID); err != nil {
			return err
		}
		// Списываем остаток
		_, err := tx.Exec(ctx, `UPDATE books SET inventory = inventory - 1 WHERE id=$1 AND inventory > 0`, item.BookID)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *OrderPostgres) ListByUser(ctx context.Context, userID string) ([]*domain.Order, error) {
	rows, err := r.db.Query(ctx, `SELECT id, created_at FROM orders WHERE user_id=$1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orders []*domain.Order
	for rows.Next() {
		var o domain.Order
		o.UserID = userID
		if err := rows.Scan(&o.ID, &o.CreatedAt); err != nil {
			return nil, err
		}
		// Получаем order_items
		itemRows, err := r.db.Query(ctx, `SELECT id, book_id, price FROM order_items WHERE order_id=$1`, o.ID)
		if err != nil {
			return nil, err
		}
		for itemRows.Next() {
			var it domain.OrderItem
			if err := itemRows.Scan(&it.ID, &it.BookID, &it.Price); err != nil {
				itemRows.Close()
				return nil, err
			}
			o.Items = append(o.Items, it)
		}
		itemRows.Close()
		orders = append(orders, &o)
	}
	return orders, nil
}
