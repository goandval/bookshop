package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourorg/bookshop/internal/domain"
)

type BookPostgres struct {
	db *pgxpool.Pool
}

func NewBookPostgres(db *pgxpool.Pool) *BookPostgres {
	return &BookPostgres{db: db}
}

func (r *BookPostgres) GetByID(ctx context.Context, id int) (*domain.Book, error) {
	row := r.db.QueryRow(ctx, `SELECT id, title, author, year, price, category_id, inventory, created_at, updated_at FROM books WHERE id=$1`, id)
	var b domain.Book
	if err := row.Scan(&b.ID, &b.Title, &b.Author, &b.Year, &b.Price, &b.CategoryID, &b.Inventory, &b.CreatedAt, &b.UpdatedAt); err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BookPostgres) List(ctx context.Context, categoryIDs []int, limit, offset int) ([]*domain.Book, error) {
	q := `SELECT id, title, author, year, price, category_id, inventory, created_at, updated_at FROM books WHERE inventory > 0`
	args := []interface{}{}
	if len(categoryIDs) > 0 {
		q += " AND category_id = ANY($1)"
		args = append(args, categoryIDs)
	}
	q += " ORDER BY id LIMIT $2 OFFSET $3"
	args = append(args, limit, offset)
	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var books []*domain.Book
	for rows.Next() {
		var b domain.Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Year, &b.Price, &b.CategoryID, &b.Inventory, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		books = append(books, &b)
	}
	return books, nil
}

func (r *BookPostgres) Create(ctx context.Context, book *domain.Book) error {
	return r.db.QueryRow(ctx, `INSERT INTO books (title, author, year, price, category_id, inventory) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, created_at, updated_at`,
		book.Title, book.Author, book.Year, book.Price, book.CategoryID, book.Inventory,
	).Scan(&book.ID, &book.CreatedAt, &book.UpdatedAt)
}

func (r *BookPostgres) Update(ctx context.Context, book *domain.Book) error {
	_, err := r.db.Exec(ctx, `UPDATE books SET title=$1, author=$2, year=$3, price=$4, category_id=$5, updated_at=NOW() WHERE id=$6`,
		book.Title, book.Author, book.Year, book.Price, book.CategoryID, book.ID)
	return err
}

func (r *BookPostgres) Delete(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM books WHERE id=$1`, id)
	return err
}
