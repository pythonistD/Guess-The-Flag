package storage

import (
	"github.com/pythonistD/Guess-The-Flag/internal/db/models"
	"sync"
)

type CountryStorage interface {
	LoadCountries([]models.Country) error
	GetRandom() (*models.Country, error)
}

type InMemoryCountryStorage struct {
	countries []models.Country
	mu        sync.RWMutex
}
