package repo

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pythonistD/Guess-The-Flag/internal/db/models"
	"github.com/pythonistD/Guess-The-Flag/internal/db/queries"
	"time"
)

type GamesRepo struct {
	db *sqlx.DB
}

func NewGamesRepo(db *sqlx.DB) *GamesRepo {
	return &GamesRepo{db: db}
}

func (g *GamesRepo) Start(ctx context.Context, game *models.Game) error {
	_, err := g.db.NamedExecContext(ctx, queries.GameQueries.Start, game)
	if err != nil {
		return fmt.Errorf("error starting game: %w", err)
	}
	return nil
}

func (g *GamesRepo) End(ctx context.Context, gameId uuid.UUID) error {
	endGame := struct {
		EndedAt time.Time `db:"ended_at"`
		GameId  uuid.UUID `db:"game_id"`
	}{
		EndedAt: time.Now().UTC(),
		GameId:  gameId,
	}
	_, err := g.db.NamedExecContext(ctx, queries.GameQueries.End, &endGame)
	if err != nil {
		return fmt.Errorf("error ending game: %w", err)
	}
	return nil
}

func (g *GamesRepo) GetUserLastGame(ctx context.Context, id uuid.UUID) (*models.Game, error) {
	var game models.Game
	err := g.db.GetContext(ctx, &game, queries.GameQueries.GetLastUserGame, id)
	if err != nil {
		return nil, fmt.Errorf("error getting last user game: %w", err)
	}
	return &game, nil
}

func (g *GamesRepo) GetById(ctx context.Context, id uuid.UUID) (*models.Game, error) {
	var game models.Game
	err := g.db.GetContext(ctx, &game, queries.GameQueries.GetByID, id)
	if err != nil {
		return nil, fmt.Errorf("error getting game: %w", err)
	}
	return &game, nil
}
