package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Uuid      uuid.UUID `db:"user_id"`
	Name      string    `db:"username"`
	Email     string    `db:"email"`
	Password  string    `db:"password_hash"`
	CreatedAt time.Time `db:"created_at"`
}

type Game struct {
	GameId       uuid.UUID `db:"game_id"`
	UserId       uuid.UUID `db:"user_id"`
	LanguageCode string    `db:"language_code"`
	StartedAt    time.Time `db:"started_at"`
	EndedAt      time.Time `db:"ended_at"`
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
	AnswerId   uuid.UUID `db:"answer_id"`
	QuestionId uuid.UUID `db:"question_id"`
	AnswerText string    `db:"answer"`
	AnsweredAt time.Time `db:"answered_at"`
	IsCorrect  bool      `db:"is_correct"`
}

type QuestionWithAnswers struct {
	AnswerText  string `db:"answer" json:"answer"`
	IsCorrect   bool   `db:"is_correct" json:"is_correct"`
	CountryName string `db:"name" json:"name"`
	CountryCode string `db:"code" json:"code"`
	FlagSVG     string `db:"flag_svg" json:"flag_svg"`
}

type Country struct {
	CountryId   int    `db:"country_id"`
	Code        string `db:"code"`
	FlagImageId int    `db:"flag_image_id"`
}

type FlagImage struct {
	ImageId   int       `db:"image_id"`
	SvgData   string    `db:"svg_data"`
	ImageHash string    `db:"image_hash"`
	FileSize  int       `db:"file_size"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type CountryNames struct {
	CountryNamesId int    `db:"country_names_id"`
	LanguageCode   string `db:"language_code"`
	CountryId      int    `db:"country_id"`
	Name           string `db:"name"`
	NormalizedName string `db:"normalized_name"`
	Threshold      int    `db:"threshold"`
	IsDisplayName  bool   `db:"is_display_name"`
}

type CountryWithImage struct {
	CountryId   int    `db:"country_id"`
	Code        string `db:"code"`
	FlagImageId int    `db:"flag_image_id"`
	SvgData     string `db:"svg_data"`
	ImageHash   string `db:"image_hash"`
}
