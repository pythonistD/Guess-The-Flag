package schema

import "github.com/google/uuid"

type QuestionReq struct {
	GameId uuid.UUID `json:"gameId"`
}

type QuestionResp struct {
	QuestionText   string    `json:"question_text"`
	FlagSVG        string    `json:"flag_svg"`
	QuestionID     uuid.UUID `json:"question_id"`
	Variant        string    `json:"variant"`
	Step           int       `json:"step"`
	MaxFlags       int       `json:"max_flags"`
	StepsCompleted int       `json:"steps_completed"`
}

type AnswerOption struct {
	CountryId int    `json:"country_id"`
	Name      string `json:"name"`
}

type MultipleChoiceQuestionResp struct {
	QuestionResp
	Options []AnswerOption `json:"options"`
}

type AnswerReq struct {
	Answer          string `json:"answer"`
	SelectedCountry int    `json:"selected_country"`
	Skipped         bool   `json:"skipped"`
}

type AnswerResp struct {
	IsCorrect      bool `json:"is_correct"`
	Step           int  `json:"step"`
	MaxFlags       int  `json:"max_flags"`
	StepsCompleted int  `json:"steps_completed"`
	GameCompleted  bool `json:"game_completed"`
}

type StartGameResp struct {
	GameId    string `json:"game_id"`
	Variant   string `json:"variant"`
	MaxFlags  int    `json:"max_flags"`
}

type EndGameReq struct {
	GameId uuid.UUID `json:"gameId"`
}

type EndGameResp struct {
}
