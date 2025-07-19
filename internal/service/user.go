package service

import (
	"context"

	"github.com/yourorg/bookshop/internal/domain"
	"github.com/yourorg/bookshop/internal/repository"
)

type UserServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{repo: repo}
}

func (s *UserServiceImpl) GetOrCreate(ctx context.Context, id, email string, isAdmin bool) (*domain.User, error) {
	return nil, nil
}
