package repo

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pythonistD/Guess-The-Flag/internal/db/models"
	"github.com/pythonistD/Guess-The-Flag/internal/db/queries"
)

type AnswersRepo struct {
	db *sqlx.DB
}

func NewAnswersRepo(db *sqlx.DB) *AnswersRepo {
	return &AnswersRepo{db: db}
}

func (a *AnswersRepo) Create(ctx context.Context, country *models.Country) error {
	_, err := a.db.ExecContext(ctx, queries.AnswerQueries.Create, country)
	if err != nil {
		return fmt.Errorf("error creating answer: %w", err)
	}
	return nil
}

func (a *AnswersRepo) GetQuestionAnswer(ctx context.Context, id uuid.UUID) (*models.Answer, error) {
	var answer models.Answer
	err := a.db.GetContext(ctx, &answer, queries.AnswerQueries.GetQuestionAnswer, id)
	if err != nil {
		return nil, fmt.Errorf("error getting answer: %w", err)
	}
	return &answer, nil
}
