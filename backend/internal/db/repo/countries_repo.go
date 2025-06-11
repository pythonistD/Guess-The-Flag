package repo

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pythonistD/Guess-The-Flag/internal/db/models"
	"github.com/pythonistD/Guess-The-Flag/internal/db/queries"
)

type CountriesRepo struct {
	db *sqlx.DB
}

func NewCountriesRepo(db *sqlx.DB) *CountriesRepo {
	return &CountriesRepo{db: db}
}

func (c *CountriesRepo) GetALl(ctx context.Context) ([]models.Country, error) {
	var countries []models.Country
	err := c.db.SelectContext(ctx, &countries, queries.CountryQueries.GetAll)
	if err != nil {
		return nil, fmt.Errorf("error getting all countries: %w", err)
	}
	return countries, nil
}
func (c *CountriesRepo) GetById(ctx context.Context, id int) (*models.Country, error) {
	var country models.Country
	err := c.db.GetContext(ctx, &country, queries.CountryQueries.GetByID, id)
	if err != nil {
		return nil, fmt.Errorf("error getting country: %w", err)
	}
	return &country, nil
}
