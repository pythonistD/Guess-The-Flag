package repo

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pythonistD/Guess-The-Flag/internal/db/models"
	"github.com/pythonistD/Guess-The-Flag/internal/db/queries"
)

type UnknownFlagsRepo struct {
	db *sqlx.DB
}

func NewUnknownFlagsRepo(db *sqlx.DB) *UnknownFlagsRepo {
	return &UnknownFlagsRepo{db: db}
}

func (u *UnknownFlagsRepo) Create(ctx context.Context, flags *models.UnknownFlags) error {
	_, err := u.db.NamedExecContext(ctx, queries.UnknownFlagQueries.Create, flags)
	if err != nil {
		return fmt.Errorf("failed to create unknown flag: %w", err)
	}
	return nil
}

func (u *UnknownFlagsRepo) GetById(ctx context.Context, id int) (*models.UnknownFlags, error) {
	var flags models.UnknownFlags
	err := u.db.GetContext(ctx, &flags, queries.UnknownFlagQueries.GetById, id)
	if err != nil {
		return nil, fmt.Errorf("error getting unknown flag: %w", err)
	}
	return &flags, nil
}

func (u *UnknownFlagsRepo) GetUserCountries(ctx context.Context, userId uuid.UUID) (*models.UnknownFlags, error) {
	var flags models.UnknownFlags
	err := u.db.GetContext(ctx, &flags, queries.UnknownFlagQueries.GetUserCountries, userId)
	if err != nil {
		return nil, fmt.Errorf("error getting user countries: %w", err)
	}
	return &flags, nil
}

func (u *UnknownFlagsRepo) Delete(ctx context.Context, id int) error {
	_, err := u.db.ExecContext(ctx, queries.UnknownFlagQueries.Delete, id)
	if err != nil {
		return fmt.Errorf("error deleting unknown flag: %w", err)
	}
	return nil
}
