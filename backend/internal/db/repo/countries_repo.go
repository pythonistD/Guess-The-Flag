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

func (c *CountriesRepo) Create(ctx context.Context, country *models.Country) (*models.Country, error) {
	rows, err := c.db.NamedQueryContext(ctx, queries.CountryQueries.Create, country)
	if err != nil {
		return nil, fmt.Errorf("error creating country: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("error creating country: no id returned")
	}

	if err := rows.Scan(&country.CountryId); err != nil {
		return nil, fmt.Errorf("error scanning created country id: %w", err)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error creating country: %w", err)
	}
	return country, nil
}

func (c *CountriesRepo) CreateAll(ctx context.Context, countries []models.Country) error {
	// Используем транзакцию для батчевой вставки
	tx, err := c.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for _, country := range countries {
		_, err := tx.NamedExecContext(ctx, queries.CountryQueries.Create, country)
		if err != nil {
			return fmt.Errorf("error creating country: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetAllWithImages получает все страны с данными изображений
func (c *CountriesRepo) GetAllWithImages(ctx context.Context) ([]models.CountryWithImage, error) {
	var countries []models.CountryWithImage
	err := c.db.SelectContext(ctx, &countries, queries.CountryQueries.GetAllWithImage)
	if err != nil {
		return nil, fmt.Errorf("error getting all countries with images: %w", err)
	}
	return countries, nil
}

// GetByIdWithImage получает страну с данными изображения
func (c *CountriesRepo) GetByIdWithImage(ctx context.Context, id int) (*models.CountryWithImage, error) {
	var country models.CountryWithImage
	err := c.db.GetContext(ctx, &country, queries.CountryQueries.GetByIdWithImage, id)
	if err != nil {
		return nil, fmt.Errorf("error getting country with image: %w", err)
	}
	return &country, nil
}
