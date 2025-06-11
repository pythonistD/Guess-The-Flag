package repo

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pythonistD/Guess-The-Flag/internal/db/models"
	"github.com/pythonistD/Guess-The-Flag/internal/db/queries"
)

type LeaderboardRepo struct {
	db *sqlx.DB
}

func NewLeaderboardRepo(db *sqlx.DB) *LeaderboardRepo {
	return &LeaderboardRepo{db: db}
}

func (l *LeaderboardRepo) UpdateOrCreate(ctx context.Context, leaderboard *models.Leaderboard) error {
	_, err := l.db.NamedExecContext(ctx, queries.LeaderboardQueries.Upsert, leaderboard)
	if err != nil {
		return fmt.Errorf("error upsert leaderboard: %w", err)
	}
	return nil
}
