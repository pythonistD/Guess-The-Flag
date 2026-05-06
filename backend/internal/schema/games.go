package schema

import "github.com/google/uuid"

type QuestionReq struct {
	GameId uuid.UUID `json:"gameId"`
}
type QuestionResp struct {
	QuestionText string    `json:"question_text"`
	FlagSVG      string    `json:"flag_svg"`
	QuestionID   uuid.UUID `json:"question_id"`
}

type AnswerReq struct {
	Answer string `json:"answer"`
}

type AnswerResp struct {
	IsCorrect bool `json:"is_correct"`
}

type StartGameResp struct {
	GameId string `json:"game_id"`
}

type EndGameReq struct {
	GameId uuid.UUID `json:"gameId"`
}

type EndGameResp struct {
}
