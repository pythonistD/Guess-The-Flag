package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	Uuid      uuid.UUID `db:"user_id"`
	Name      string    `db:"username"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
}

type Game struct {
	GameId    uuid.UUID `db:"game_id"`
	UserId    uuid.UUID `db:"user_id"`
	StartedAt time.Time `db:"started_at"`
	EndedAt   time.Time `db:"ended_at"`
}

type Leaderboard struct {
	LeaderboardId int       `db:"leaderboard_id"`
	UserId        uuid.UUID `db:"user_id"`
	Score         int       `db:"score"`
	GamesPlayed   int       `db:"games_played"`
	LastGame      time.Time `db:"last_game"`
}

type UnknownFlags struct {
	UnknownFlagsId int       `db:"unknown_flags_id"`
	UserId         uuid.UUID `db:"user_id"`
	CountryId      int       `db:"country_id"`
}

type Question struct {
	QuestionId uuid.UUID `db:"question_id"`
	GameId     uuid.UUID `db:"game_id"`
	CountryId  int       `db:"country_id"`
	CreatedAt  time.Time `db:"created_at"`
}

type Answer struct {
	AnswerId        uuid.UUID `db:"answer_id"`
	QuestionId      uuid.UUID `db:"question_id"`
	SelectedCountry int       `db:"selected_country_id"`
	AnsweredAt      time.Time `db:"answered_at"`
	IsCorrect       bool      `db:"is_correct"`
}

type Country struct {
	CountryId int    `db:"country_id"`
	Name      string `db:"name"`
	Code      string `db:"code"`
	FlagUrl   string `db:"flag_url"`
}
