package storage

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type GameStorage interface {
	SetQuestions(gameID uuid.UUID, questions []QuestionInStorage) error
	GetQuestions(gameID uuid.UUID) ([]QuestionInStorage, error)
	GetQuestion(gameID uuid.UUID, questionNum int) (*QuestionInStorage, error)
	DeleteQuestions(gameID uuid.UUID) error
	GetQuestionsRemaining(gameID uuid.UUID) (int, error)
	// PopQuestion extracts and deletes the question
	PopQuestion(gameID uuid.UUID, questionNum int) (*QuestionInStorage, error)

	SetCountry(gameID uuid.UUID, countryId int) error
	// IsCountryUsed Checks if country already used in current game session
	IsCountryUsed(gameID uuid.UUID, countryId int) bool
	DeleteCountries(gameID uuid.UUID, countryId int) error
}

type QuestionInStorage struct {
	QuestionId   uuid.UUID
	GameId       uuid.UUID
	QuestionText string
	FlagUrl      string
	Answer       string
	CountryId    int
	CreatedAt    time.Time
}

type InMemoryGameStorage struct {
	mu            sync.RWMutex
	questions     map[uuid.UUID]map[int]QuestionInStorage
	countriesUsed map[uuid.UUID][]int
}

// NewInMemoryGameStorage creates a new instance of InMemoryGameStorage
func NewInMemoryGameStorage() *InMemoryGameStorage {
	return &InMemoryGameStorage{
		questions:     make(map[uuid.UUID]map[int]QuestionInStorage),
		countriesUsed: make(map[uuid.UUID][]int),
	}
}

func (i *InMemoryGameStorage) GetQuestionsRemaining(gameID uuid.UUID) (int, error) {
	if _, exists := i.questions[gameID]; !exists {
		return 0, nil
	}
	return len(i.questions[gameID]), nil
}

func (i *InMemoryGameStorage) SetQuestions(gameID uuid.UUID, questions []QuestionInStorage) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	// Initialize the map for this game if it doesn't exist
	if i.questions[gameID] == nil {
		i.questions[gameID] = make(map[int]QuestionInStorage)
	}

	// Store questions with their index as key
	for idx, question := range questions {
		i.questions[gameID][idx] = question
	}

	return nil
}

func (i *InMemoryGameStorage) GetQuestions(gameID uuid.UUID) ([]QuestionInStorage, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	gameQuestions, exists := i.questions[gameID]
	if !exists {
		return nil, fmt.Errorf("no questions found for game ID: %s", gameID)
	}

	// Convert map to slice, maintaining order
	questions := make([]QuestionInStorage, len(gameQuestions))
	for idx, question := range gameQuestions {
		if idx < len(questions) {
			questions[idx] = question
		}
	}

	return questions, nil
}

var (
	GameIdError     = errors.New("no questions found for game ID")
	QuestionIdError = errors.New("no questions with this questionNum")
)

func (i *InMemoryGameStorage) GetQuestion(gameID uuid.UUID, questionNum int) (*QuestionInStorage, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	gameQuestions, exists := i.questions[gameID]
	if !exists {
		return nil, GameIdError
	}

	question, exists := gameQuestions[questionNum]
	if !exists {
		return nil, QuestionIdError
	}

	return &question, nil
}

func (i *InMemoryGameStorage) DeleteQuestions(gameID uuid.UUID) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	delete(i.questions, gameID)
	return nil
}

func (i *InMemoryGameStorage) PopQuestion(gameID uuid.UUID, questionNum int) (*QuestionInStorage, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	gameQuestions, exists := i.questions[gameID]
	if !exists {
		return nil, GameIdError
	}

	question, exists := gameQuestions[questionNum]
	if !exists {
		return nil, QuestionIdError
	}

	// Remove the question from storage
	delete(gameQuestions, questionNum)

	// If this was the last question, clean up the game entry
	if len(gameQuestions) == 0 {
		delete(i.questions, gameID)
	}

	return &question, nil
}

func (i *InMemoryGameStorage) SetCountry(gameID uuid.UUID, countryId int) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	// Initialize the slice for this game if it doesn't exist
	if i.countriesUsed[gameID] == nil {
		i.countriesUsed[gameID] = make([]int, 0)
	}

	// Check if country is already used
	for _, usedCountryId := range i.countriesUsed[gameID] {
		if usedCountryId == countryId {
			return fmt.Errorf("country ID %d is already used in game %s", countryId, gameID)
		}
	}

	// Add country to used list
	i.countriesUsed[gameID] = append(i.countriesUsed[gameID], countryId)
	return nil
}

func (i *InMemoryGameStorage) IsCountryUsed(gameID uuid.UUID, countryId int) bool {
	i.mu.RLock()
	defer i.mu.RUnlock()

	usedCountries, exists := i.countriesUsed[gameID]
	if !exists {
		return false
	}

	for _, usedCountryId := range usedCountries {
		if usedCountryId == countryId {
			return true
		}
	}

	return false
}

func (i *InMemoryGameStorage) DeleteCountries(gameID uuid.UUID, countryId int) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	usedCountries, exists := i.countriesUsed[gameID]
	if !exists {
		return fmt.Errorf("no countries found for game ID: %s", gameID)
	}

	// Find and remove the specific country
	found := false
	newCountries := make([]int, 0, len(usedCountries))
	for _, usedCountryId := range usedCountries {
		if usedCountryId == countryId {
			found = true
		} else {
			newCountries = append(newCountries, usedCountryId)
		}
	}

	if !found {
		return fmt.Errorf("country ID %d not found in game %s", countryId, gameID)
	}

	i.countriesUsed[gameID] = newCountries
	return nil
}
