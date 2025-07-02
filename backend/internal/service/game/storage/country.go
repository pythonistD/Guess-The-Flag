package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pythonistD/Guess-The-Flag/internal/db/models"
	"github.com/pythonistD/Guess-The-Flag/internal/db/repo"
	"math/rand"
	"sync"
)

type CountryStorage interface {
	LoadCountriesFromDB() error
	loadCountries([]models.Country) error
	GetRandom() (*models.Country, error)
}

type InMemoryCountryStorage struct {
	countriesRepo *repo.CountriesRepo
	countries     []models.Country
	mu            sync.RWMutex
}

func NewInMemoryCountryStorage(db *sqlx.DB) *InMemoryCountryStorage {
	return &InMemoryCountryStorage{
		countriesRepo: repo.NewCountriesRepo(db),
		countries:     make([]models.Country, 0, 200),
	}
}

func (i *InMemoryCountryStorage) LoadCountriesFromDB() error {
	ctx := context.Background()
	countries, err := i.countriesRepo.GetALl(ctx)
	if err != nil {
		return fmt.Errorf("error while getting countries from: %w", err)
	}
	err = i.loadCountries(countries)
	if err != nil {
		return fmt.Errorf("error while loading countries to inmemory storage: %w", err)
	}
	return nil
}

func (i *InMemoryCountryStorage) loadCountries(countries []models.Country) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	for _, c := range countries {
		i.countries = append(i.countries, c)
	}
	return nil
}

func (i *InMemoryCountryStorage) GetRandom() (*models.Country, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	if len(i.countries) == 0 {
		return nil, errors.New("no countries found. You have to load countries first")
	}
	countryNum := rand.Intn(len(i.countries))
	return &i.countries[countryNum], nil
}
