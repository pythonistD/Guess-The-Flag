package repo

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pythonistD/Guess-The-Flag/internal/db/models"
	"github.com/pythonistD/Guess-The-Flag/internal/db/queries"
)

type QuestionsRepo struct {
	db *sqlx.DB
}

func NewQuestionsRepo(db *sqlx.DB) *QuestionsRepo {
	return &QuestionsRepo{db: db}
}

func (q *QuestionsRepo) Create(ctx context.Context, question *models.Question) error {
	_, err := q.db.NamedQueryContext(ctx, queries.QuestionQueries.Create, question)
	if err != nil {
		return fmt.Errorf("failed to create question: %w", err)
	}
	return nil
}

func (q *QuestionsRepo) GetGameQuestions(ctx context.Context, id uuid.UUID) (*models.Question, error) {
	var question models.Question
	err := q.db.GetContext(ctx, &question, queries.QuestionQueries.GetGameQuestions, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get question: %w", err)
	}
	return &question, nil
}

func (q *QuestionsRepo) GetQuestionsWithAnswers(ctx context.Context, gameId uuid.UUID, langCode string) ([]models.QuestionWithAnswers, error) {
	var questionsWithAnswer []models.QuestionWithAnswers
	err := q.db.SelectContext(ctx, &questionsWithAnswer, queries.QuestionQueries.GetQuestionsWithAnswers, gameId, langCode)
	if err != nil {
		return nil, fmt.Errorf("error getting questions with answers: %w", err)
	}
	return questionsWithAnswer, nil
}

func (q *QuestionsRepo) GetQuestion(ctx context.Context, questionId uuid.UUID) (*models.Question, error) {
	var question models.Question
	err := q.db.GetContext(ctx, &question, queries.QuestionQueries.GetQuestion, questionId)
	if err != nil {
		return nil, fmt.Errorf("failed to get question: %w", err)
	}
	return &question, nil
}
