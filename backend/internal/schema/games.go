package schema

import "github.com/google/uuid"

type QuestionReq struct {
	GameId      uuid.UUID `json:"gameId"`
	QuestionNum int       `json:"questionNum"`
}
type QuestionResp struct {
	QuestionText string `json:"question_text"`
	FlagUrl      string `json:"flag_url"`
}

type AnswerReq struct {
	GameId      uuid.UUID `json:"gameId"`
	QuestionNum int       `json:"questionNum"`
	Answer      string    `json:"answer"`
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
