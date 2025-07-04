package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourorg/bookshop/internal/domain"
)

type CategoryPostgres struct {
	db *pgxpool.Pool
}

func NewCategoryPostgres(db *pgxpool.Pool) *CategoryPostgres {
	return &CategoryPostgres{db: db}
}

func (r *CategoryPostgres) GetByID(ctx context.Context, id int) (*domain.Category, error) {
	row := r.db.QueryRow(ctx, `SELECT id, name FROM categories WHERE id=$1`, id)
	var c domain.Category
	if err := row.Scan(&c.ID, &c.Name); err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoryPostgres) List(ctx context.Context) ([]*domain.Category, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name FROM categories ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var cats []*domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		cats = append(cats, &c)
	}
	return cats, nil
}

func (r *CategoryPostgres) Create(ctx context.Context, category *domain.Category) error {
	return r.db.QueryRow(ctx, `INSERT INTO categories (name) VALUES ($1) RETURNING id`, category.Name).Scan(&category.ID)
}

func (r *CategoryPostgres) Update(ctx context.Context, category *domain.Category) error {
	_, err := r.db.Exec(ctx, `UPDATE categories SET name=$1 WHERE id=$2`, category.Name, category.ID)
	return err
}

func (r *CategoryPostgres) Delete(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM categories WHERE id=$1`, id)
	return err
}

func (r *CategoryPostgres) GetByName(ctx context.Context, name string) (*domain.Category, error) {
	row := r.db.QueryRow(ctx, `SELECT id, name FROM categories WHERE name=$1`, name)
	var c domain.Category
	if err := row.Scan(&c.ID, &c.Name); err != nil {
		return nil, err
	}
	return &c, nil
}
