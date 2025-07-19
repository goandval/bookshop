package repository

import (
	"context"
	"fmt"
	"strconv"

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
		return nil, fmt.Errorf("get by id: %w", err)
	}
	return &b, nil
}

func (r *BookPostgres) List(ctx context.Context, categoryIDs []int, limit, offset int) ([]*domain.Book, error) {
	q := `SELECT id, title, author, year, price, category_id, inventory, created_at, updated_at FROM books WHERE inventory > 0`
	args := []interface{}{}
	paramCount := 0

	if len(categoryIDs) > 0 {
		paramCount++
		q += " AND category_id = ANY($" + strconv.Itoa(paramCount) + ")"
		args = append(args, categoryIDs)
	}

	paramCount++
	q += " ORDER BY id LIMIT $" + strconv.Itoa(paramCount)
	args = append(args, limit)

	paramCount++
	q += " OFFSET $" + strconv.Itoa(paramCount)
	args = append(args, offset)

	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("list books: %w", err)
	}
	defer rows.Close()
	var books []*domain.Book
	for rows.Next() {
		var b domain.Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Year, &b.Price, &b.CategoryID, &b.Inventory, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan book: %w", err)
		}
		books = append(books, &b)
	}
	if books == nil {
		books = make([]*domain.Book, 0)
	}
	return books, nil
}

func (r *BookPostgres) Create(ctx context.Context, book *domain.Book) error {
	err := r.db.QueryRow(ctx, `INSERT INTO books (title, author, year, price, category_id, inventory) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, created_at, updated_at`,
		book.Title, book.Author, book.Year, book.Price, book.CategoryID, book.Inventory,
	).Scan(&book.ID, &book.CreatedAt, &book.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create book: %w", err)
	}
	return nil
}

func (r *BookPostgres) Update(ctx context.Context, book *domain.Book) error {
	_, err := r.db.Exec(ctx, `UPDATE books SET title=$1, author=$2, year=$3, price=$4, category_id=$5, updated_at=NOW() WHERE id=$6`,
		book.Title, book.Author, book.Year, book.Price, book.CategoryID, book.ID)
	if err != nil {
		return fmt.Errorf("update book: %w", err)
	}
	return nil
}

func (r *BookPostgres) Delete(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM books WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("delete book: %w", err)
	}
	return nil
}
