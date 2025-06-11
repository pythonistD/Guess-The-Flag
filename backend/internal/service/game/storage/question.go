package storage

import (
	"github.com/google/uuid"
	"github.com/pythonistD/Guess-The-Flag/internal/db/models"
	"sync"
)

type QuestionStorage interface {
	SetQuestions(gameID uuid.UUID, questions []models.Question) error
	GetQuestions(gameID uuid.UUID) ([]models.Question, error)
	GetQuestion(gameID uuid.UUID, questionPos int) (models.Question, error)
	DeleteQuestions(gameID uuid.UUID) error
}

type InMemoryQuestionStorage struct {
	mu    sync.RWMutex
	store map[uuid.UUID][]models.Question
}
