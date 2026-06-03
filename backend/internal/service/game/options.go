package game

import (
	"fmt"
	"math/rand"

	"github.com/pythonistD/Guess-The-Flag/internal/schema"
	"github.com/pythonistD/Guess-The-Flag/internal/service/game/storage"
)

const distractorCount = 3

func GenerateOptions(
	store storage.CountryStorage,
	correctCountryId int,
	langCode string,
	usedCountryIds map[int]struct{},
) ([]schema.AnswerOption, error) {
	correctCountry, err := store.GetByID(correctCountryId)
	if err != nil {
		return nil, fmt.Errorf("failed to get correct country: %w", err)
	}
	commonName, ok := correctCountry.CommonCountryNames[langCode]
	if !ok {
		return nil, fmt.Errorf("no common name for country %d in lang %s", correctCountryId, langCode)
	}

	options := make([]schema.AnswerOption, 0, distractorCount+1)
	options = append(options, schema.AnswerOption{
		CountryId: correctCountryId,
		Name:      commonName.Name,
	})

	excluded := make(map[int]struct{}, len(usedCountryIds)+1)
	excluded[correctCountryId] = struct{}{}
	for id := range usedCountryIds {
		excluded[id] = struct{}{}
	}

	distractors := pickDistractors(store, excluded, distractorCount, langCode)
	options = append(options, distractors...)

	rand.Shuffle(len(options), func(i, j int) {
		options[i], options[j] = options[j], options[i]
	})

	return options, nil
}

func pickDistractors(
	store storage.CountryStorage,
	excluded map[int]struct{},
	count int,
	langCode string,
) []schema.AnswerOption {
	distractors := make([]schema.AnswerOption, 0, count)
	attempts := 0
	maxAttempts := storage.NumberOfCountries * 2

	for len(distractors) < count && attempts < maxAttempts {
		attempts++
		country, err := store.GetRandom()
		if err != nil {
			continue
		}
		if _, ok := excluded[country.Id]; ok {
			continue
		}
		excluded[country.Id] = struct{}{}

		name := countryNameForLang(country, langCode)
		if name == "" {
			continue
		}
		distractors = append(distractors, schema.AnswerOption{
			CountryId: country.Id,
			Name:      name,
		})
	}

	if len(distractors) < count {
		relaxExcluded := make(map[int]struct{})
		for id := range excluded {
			relaxExcluded[id] = struct{}{}
		}
		for len(distractors) < count && attempts < maxAttempts*2 {
			attempts++
			country, err := store.GetRandom()
			if err != nil {
				continue
			}
			if _, ok := relaxExcluded[country.Id]; ok {
				continue
			}
			relaxExcluded[country.Id] = struct{}{}
			name := countryNameForLang(country, langCode)
			if name == "" {
				continue
			}
			distractors = append(distractors, schema.AnswerOption{
				CountryId: country.Id,
				Name:      name,
			})
		}
	}

	return distractors
}

func countryNameForLang(country storage.Country, langCode string) string {
	if commonName, ok := country.CommonCountryNames[langCode]; ok {
		return commonName.Name
	}
	for _, name := range country.CommonCountryNames {
		return name.Name
	}
	return ""
}
