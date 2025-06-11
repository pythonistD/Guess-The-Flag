package repo

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pythonistD/Guess-The-Flag/internal/db/models"
	"github.com/pythonistD/Guess-The-Flag/internal/db/queries"
)

type UsersRepo struct {
	db *sqlx.DB
}

func NewUsersRepo(db *sqlx.DB) *UsersRepo {
	return &UsersRepo{db: db}
}

func (u *UsersRepo) Create(ctx context.Context, user *models.User) error {
	_, err := u.db.NamedQueryContext(ctx, queries.UserQueries.Create, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (u *UsersRepo) GetById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := u.db.GetContext(ctx, &user, queries.UserQueries.GetByID, id)
	if err != nil {
		return nil, fmt.Errorf("error getting user by id: %w", err)
	}
	return &user, nil
}

func (u *UsersRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := u.db.GetContext(ctx, &user, queries.UserQueries.GetByEmail, email)
	if err != nil {
		return nil, fmt.Errorf("error getting user by email: %w", err)
	}
	return &user, nil
}
