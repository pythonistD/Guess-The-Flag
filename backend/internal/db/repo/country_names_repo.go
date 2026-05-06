package repo

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pythonistD/Guess-The-Flag/internal/db/models"
	"github.com/pythonistD/Guess-The-Flag/internal/db/queries"
)

type CountryNamesRepo struct {
	db *sqlx.DB
}

func NewCountryNamesRepo(db *sqlx.DB) *CountryNamesRepo {
	return &CountryNamesRepo{
		db: db,
	}
}

func (c *CountryNamesRepo) CreateAll(ctx context.Context, countryNames []models.CountryNames) error {
	tx, err := c.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for _, countryName := range countryNames {
		_, err := tx.NamedExecContext(ctx, queries.CountryNamesQueries.Create, countryName)
		if err != nil {
			return fmt.Errorf("error creating country name: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (c *CountryNamesRepo) GetCommonNamesByCountryId(ctx context.Context, countryId int) ([]models.CountryNames, error) {
	var names []models.CountryNames
	err := c.db.SelectContext(ctx, &names, queries.CountryNamesQueries.GetCommonNamesByCountryId, countryId)
	if err != nil {
		return nil, fmt.Errorf("error getting common names by country id: %w", err)
	}
	return names, nil
}

func (c *CountryNamesRepo) GetAllNamesByCountryId(ctx context.Context, countryId int) ([]models.CountryNames, error) {
	var names []models.CountryNames
	err := c.db.SelectContext(ctx, &names, queries.CountryNamesQueries.GetAllNamesByCountryId, countryId)
	if err != nil {
		return nil, fmt.Errorf("error getting all names by country id: %w", err)
	}
	return names, nil
}
