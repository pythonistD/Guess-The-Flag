package storage

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

const QuestionsBuffSize = 5

type GameStorage interface {
	InitStorageState(gameId uuid.UUID, langCode string)
	SetQuestions(gameID uuid.UUID, questions []QuestionInStorage) error
	GetQuestion(gameID uuid.UUID) (*QuestionInStorage, error)
	GetQuestionsRemaining(gameID uuid.UUID) (int, error)
	PopQuestion(gameID uuid.UUID) (*QuestionInStorage, error)

	SetCountry(gameID uuid.UUID, countryId int) error
	// IsCountryUsed Checks if country already used in current game session
	IsCountryUsed(gameID uuid.UUID, countryId int) bool

	DeleteGameSession(gameID uuid.UUID) error

	GetGameLangCode(gameID uuid.UUID) (string, error)
}

type QuestionInStorage struct {
	QuestionId   uuid.UUID
	GameId       uuid.UUID
	QuestionText string
	FlagSVG      string
	Answer       string // название страны, которое отобразить
	CountryId    int
	CreatedAt    time.Time
}

type GameSession struct {
	// questions - кольцевой буфер
	questions     *RingBuffer[QuestionInStorage]
	countriesUsed map[int]struct{}
	LangCode      string
}

type InMemoryGameStorage struct {
	mu           sync.RWMutex
	GameSessions map[uuid.UUID]GameSession
}

// NewInMemoryGameStorage creates a new instance of InMemoryGameStorage
func NewInMemoryGameStorage() *InMemoryGameStorage {
	return &InMemoryGameStorage{
		GameSessions: make(map[uuid.UUID]GameSession),
	}
}

func (i *InMemoryGameStorage) InitStorageState(gameId uuid.UUID, langCode string) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.GameSessions[gameId] = GameSession{
		questions:     NewRingBuffer[QuestionInStorage](QuestionsBuffSize + 1),
		countriesUsed: make(map[int]struct{}, NumberOfCountries),
		LangCode:      langCode,
	}
}

func (i *InMemoryGameStorage) GetQuestionsRemaining(gameID uuid.UUID) (int, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.GameSessions[gameID].questions.RemainingItemsNumber(), nil
}

func (i *InMemoryGameStorage) SetQuestions(gameID uuid.UUID, questions []QuestionInStorage) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	session, exists := i.GameSessions[gameID]
	if !exists {
		session = GameSession{
			questions:     NewRingBuffer[QuestionInStorage](QuestionsBuffSize + 1),
			countriesUsed: make(map[int]struct{}, NumberOfCountries),
		}
	} else if session.questions == nil {
		session.questions = NewRingBuffer[QuestionInStorage](QuestionsBuffSize + 1)
	}

	for _, question := range questions {
		err := session.questions.Push(question)
		if err != nil {
			return fmt.Errorf("failed to push question: %w", err)
		}
	}
	i.GameSessions[gameID] = session
	return nil
}

var (
	GameIdError           = errors.New("no questions found for game ID")
	GettingQuestionsError = errors.New("failed to get questions")
)

func (i *InMemoryGameStorage) GetQuestion(gameID uuid.UUID) (*QuestionInStorage, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	session, exists := i.GameSessions[gameID]
	if !exists {
		return nil, GameIdError
	}
	q, err := session.questions.Pop()
	if err != nil {
		return nil, errors.Join(GettingQuestionsError, err)
	}
	return &q, nil
}

func (i *InMemoryGameStorage) DeleteQuestions(gameID uuid.UUID) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	delete(i.GameSessions, gameID)
	return nil
}

func (i *InMemoryGameStorage) PopQuestion(gameID uuid.UUID) (*QuestionInStorage, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	session, exists := i.GameSessions[gameID]
	if !exists {
		return nil, GameIdError
	}
	gameQuestions := session.questions
	q, err := gameQuestions.Pop()
	if err != nil {
		return nil, err
	}
	return &q, nil
}

func (i *InMemoryGameStorage) SetCountry(gameID uuid.UUID, countryId int) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	session, exists := i.GameSessions[gameID]
	if !exists {
		session = GameSession{
			questions:     NewRingBuffer[QuestionInStorage](QuestionsBuffSize + 1),
			countriesUsed: make(map[int]struct{}, NumberOfCountries),
		}
	} else if session.countriesUsed == nil {
		session.countriesUsed = make(map[int]struct{}, NumberOfCountries)
	}
	if _, ok := session.countriesUsed[countryId]; ok {
		return fmt.Errorf("country %d already used", countryId)
	}
	// Add country to used list
	session.countriesUsed[countryId] = struct{}{}
	i.GameSessions[gameID] = session
	return nil
}

func (i *InMemoryGameStorage) IsCountryUsed(gameID uuid.UUID, countryId int) bool {
	i.mu.RLock()
	defer i.mu.RUnlock()

	session, exists := i.GameSessions[gameID]
	if !exists {
		return false
	}
	usedCountries := session.countriesUsed
	if _, ok := usedCountries[countryId]; ok {
		return true
	}
	return false
}
func (i *InMemoryGameStorage) DeleteGameSession(gameID uuid.UUID) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	delete(i.GameSessions, gameID)
	return nil
}

func (i *InMemoryGameStorage) GetGameLangCode(gameID uuid.UUID) (string, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	session, exists := i.GameSessions[gameID]
	if !exists {
		return "", GameIdError
	}
	return session.LangCode, nil
}
