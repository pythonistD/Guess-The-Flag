package storage

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/pythonistD/Guess-The-Flag/internal/db/repo"
)

const NumberOfCountries = 195

type CountryStorage interface {
	InitCountryStorageState() error
	GetRandom() (Country, error)
	GetByID(id int) (Country, error)
	GetAll() []Country
}

type CountryName struct {
	Name           string
	NormalizedName string
	Threshold      int
}

type Country struct {
	// В обоих словарях ключ - код языка из 3-х символов
	CommonCountryNames map[string]CountryName
	AllCountryNames    map[string][]CountryName
	FlagSVG            string
	Id                 int
}

type InMemoryCountryStorage struct {
	countriesRepo    *repo.CountriesRepo
	countryNamesRepo *repo.CountryNamesRepo
	countries        []Country // ключ - country_id
	mu               sync.RWMutex
}

func NewInMemoryCountryStorage(db *sqlx.DB) *InMemoryCountryStorage {
	return &InMemoryCountryStorage{
		countriesRepo:    repo.NewCountriesRepo(db),
		countryNamesRepo: repo.NewCountryNamesRepo(db),
		countries:        make([]Country, 0, NumberOfCountries),
	}
}

func (i *InMemoryCountryStorage) InitCountryStorageState() error {
	ctx := context.Background()
	countriesWithImage, err := i.countriesRepo.GetAllWithImages(ctx)
	if err != nil {
		return fmt.Errorf("country storage init error: %w", err)
	}
	for _, country := range countriesWithImage {
		common, err := i.loadCommonNames(ctx, country.CountryId)
		if err != nil {
			return fmt.Errorf("country storage init error: %w", err)
		}
		all, err := i.loadAllNames(ctx, country.CountryId)
		if err != nil {
			return fmt.Errorf("country storage init error: %w", err)
		}
		inMemoryCountry := Country{
			CommonCountryNames: common,
			AllCountryNames:    all,
			FlagSVG:            country.SvgData,
			Id:                 country.CountryId,
		}
		i.countries = append(i.countries, inMemoryCountry)
	}
	return nil
}

func (i *InMemoryCountryStorage) loadCommonNames(ctx context.Context, countryId int) (map[string]CountryName, error) {
	names, err := i.countryNamesRepo.GetCommonNamesByCountryId(ctx, countryId)
	if err != nil {
		return nil, err
	}
	commonNames := make(map[string]CountryName)
	for _, name := range names {
		commonNames[name.LanguageCode] = CountryName{
			Name:           name.Name,
			NormalizedName: name.NormalizedName,
			Threshold:      name.Threshold,
		}
	}
	return commonNames, nil
}

func (i *InMemoryCountryStorage) loadAllNames(ctx context.Context, countryId int) (map[string][]CountryName, error) {
	names, err := i.countryNamesRepo.GetAllNamesByCountryId(ctx, countryId)
	if err != nil {
		return nil, err
	}
	allNames := make(map[string][]CountryName)
	for _, name := range names {
		allNames[name.LanguageCode] = append(allNames[name.LanguageCode], CountryName{
			Name:           name.Name,
			NormalizedName: name.NormalizedName,
			Threshold:      name.Threshold,
		})
	}
	return allNames, nil
}

func (i *InMemoryCountryStorage) GetRandom() (Country, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	if len(i.countries) == 0 {
		return Country{
			CommonCountryNames: make(map[string]CountryName),
			AllCountryNames:    make(map[string][]CountryName),
			FlagSVG:            "",
		}, errors.New("no countries found. You have to load countries first")
	}
	n := rand.Intn(len(i.countries))
	return i.countries[n], nil
}

// GetAll возвращает копию слайса со всеми странами из in-memory хранилища.
// Используется для отладочных эндпоинтов (например, отрисовки всех флагов).
func (i *InMemoryCountryStorage) GetAll() []Country {
	i.mu.RLock()
	defer i.mu.RUnlock()
	result := make([]Country, len(i.countries))
	copy(result, i.countries)
	return result
}

func (i *InMemoryCountryStorage) GetByID(id int) (Country, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	if len(i.countries) == 0 {
		return Country{}, errors.New("no countries found. You have to load countries first")
	}
	for _, c := range i.countries {
		if c.Id == id {
			return c, nil
		}
	}
	return Country{}, fmt.Errorf("country with id %d not found", id)
}
