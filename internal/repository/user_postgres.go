package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourorg/bookshop/internal/domain"
)

type UserPostgres struct {
	db *pgxpool.Pool
}

func NewUserPostgres(db *pgxpool.Pool) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) GetByID(ctx context.Context, id string) (*domain.User, error) {
	row := r.db.QueryRow(ctx, `SELECT id, email, is_admin FROM users WHERE id=$1`, id)
	var u domain.User
	if err := row.Scan(&u.ID, &u.Email, &u.IsAdmin); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserPostgres) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	row := r.db.QueryRow(ctx, `SELECT id, email, is_admin FROM users WHERE email=$1`, email)
	var u domain.User
	if err := row.Scan(&u.ID, &u.Email, &u.IsAdmin); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserPostgres) CreateIfNotExists(ctx context.Context, user *domain.User) error {
	_, err := r.db.Exec(ctx, `INSERT INTO users (id, email, is_admin) VALUES ($1, $2, $3) ON CONFLICT (id) DO NOTHING`, user.ID, user.Email, user.IsAdmin)
	return err
}
