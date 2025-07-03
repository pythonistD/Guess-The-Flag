package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pythonistD/Guess-The-Flag/internal/schema"
	"github.com/pythonistD/Guess-The-Flag/internal/service/game"
	"go.uber.org/zap"
)

type GameHandler struct {
	gameService *game.Service
	logger      *zap.Logger
}

func NewGameHandler(gameService *game.Service, logger *zap.Logger) *GameHandler {
	return &GameHandler{
		gameService: gameService,
		logger:      logger,
	}
}

func (handler *GameHandler) Start(w http.ResponseWriter, r *http.Request) {
	// Обрабатываем OPTIONS запросы для CORS
	if r.Method == "OPTIONS" {
		return
	}

	gameId, err := handler.gameService.StartGame(r.Context())
	if err != nil {
		handler.logger.Error("Failed to start game", zap.Error(err))
		http.Error(w, "Failed to start game", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(
		schema.StartGameResp{
			GameId: gameId.String(),
		},
	)
	if err != nil {
		handler.logger.Error("failed to start the game", zap.Error(err))
		http.Error(w, fmt.Sprintf("failed to start the game: %s", err.Error()), http.StatusInternalServerError)
	}
}

func (handler *GameHandler) GetQuestion(w http.ResponseWriter, r *http.Request) {
	// Обрабатываем OPTIONS запросы для CORS
	if r.Method == "OPTIONS" {
		return
	}

	var req schema.QuestionReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handler.logger.Error("failed to decode the request", zap.Error(err))
		http.Error(w, fmt.Sprintf("failed to decode the request: %s", err.Error()), http.StatusBadRequest)
	}
	question, err := handler.gameService.GetQuestion(r.Context(), req.GameId, req.QuestionNum)
	if err != nil {
		handler.logger.Error("failed to get the question", zap.Error(err))
		http.Error(w, fmt.Sprintf("failed to get the question: %s", err.Error()), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(question)
	if err != nil {
		handler.logger.Error("failed to encode the response", zap.Error(err))
		http.Error(w, fmt.Sprintf("failed to encode the response: %s", err.Error()), http.StatusInternalServerError)
	}
}

func (handler *GameHandler) AnswerTheQuestion(w http.ResponseWriter, r *http.Request) {
	// Обрабатываем OPTIONS запросы для CORS
	if r.Method == "OPTIONS" {
		return
	}

	var req schema.AnswerReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handler.logger.Error("failed to decode the request", zap.Error(err))
		http.Error(w, fmt.Sprintf("failed to decode the request: %s", err.Error()), http.StatusBadRequest)
	}
	resp, err := handler.gameService.AnswerTheQuestion(r.Context(), req.GameId, req.QuestionNum, req.Answer)
	if err != nil {
		handler.logger.Error("failed to answer the question", zap.Error(err))
		http.Error(w, fmt.Sprintf("failed to answer the question: %s", err.Error()), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
}

func (handler *GameHandler) End(w http.ResponseWriter, r *http.Request) {
	// Обрабатываем OPTIONS запросы для CORS
	if r.Method == "OPTIONS" {
		return
	}

	var req schema.EndGameReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handler.logger.Error("failed to decode the request", zap.Error(err))
		http.Error(w, fmt.Sprintf("failed to decode the request: %s", err.Error()), http.StatusBadRequest)
	}
	questionsWithAnswer, err := handler.gameService.EndGame(r.Context(), req.GameId)
	if err != nil {
		handler.logger.Error("failed to end the game", zap.Error(err))
		http.Error(w, fmt.Sprintf("failed to end the game: %s", err.Error()), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(questionsWithAnswer)
}
