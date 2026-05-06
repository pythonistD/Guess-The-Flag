package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

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
	if r.Method == "OPTIONS" {
		return
	}
	langCode := r.URL.Query().Get("lang_code")
	if langCode == "" {
		langCode = "rus"
	}

	gameId, err := handler.gameService.StartGame(r.Context(), langCode)
	if err != nil {
		handler.logger.Error("Failed to start game", zap.Error(err))
		http.Error(w, "Failed to start game", http.StatusInternalServerError)
		return
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
		return
	}
}

func (handler *GameHandler) GetQuestion(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

	gameIdS := mux.Vars(r)["gameId"]
	if gameIdS == "" {
		handler.logger.Error("Missing gameId in path")
		http.Error(w, "Missing gameId", http.StatusBadRequest)
		return
	}
	gameId, err := uuid.Parse(gameIdS)
	if err != nil {
		handler.logger.Error("Invalid gameId", zap.String("gameId", gameIdS))
		http.Error(w, "Invalid gameId", http.StatusBadRequest)
		return
	}

	question, err := handler.gameService.GetQuestion(r.Context(), gameId)
	if err != nil {
		handler.logger.Error("failed to get the question", zap.Error(err))
		http.Error(w, fmt.Sprintf("failed to get the question: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(question)
	if err != nil {
		handler.logger.Error("failed to encode the response", zap.Error(err))
		http.Error(w, fmt.Sprintf("failed to encode the response: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (handler *GameHandler) AnswerTheQuestion(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

	gameIdS := mux.Vars(r)["gameId"]
	if gameIdS == "" {
		handler.logger.Error("Missing gameId in path")
		http.Error(w, "Missing gameId", http.StatusBadRequest)
		return
	}
	gameId, err := uuid.Parse(gameIdS)
	if err != nil {
		handler.logger.Error("Invalid gameId", zap.String("gameId", gameIdS))
		http.Error(w, "Invalid gameId", http.StatusBadRequest)
		return
	}
	questionIdS := mux.Vars(r)["questionId"]
	if questionIdS == "" {
		handler.logger.Error("Missing questionId in path")
		http.Error(w, "Missing questionId", http.StatusBadRequest)
		return
	}
	questionId, err := uuid.Parse(questionIdS)
	if err != nil {
		handler.logger.Error("Invalid questionId", zap.String("questionId", questionIdS))
		http.Error(w, "Invalid questionId", http.StatusBadRequest)
		return
	}

	var req schema.AnswerReq
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handler.logger.Error("failed to decode the request", zap.Error(err))
		http.Error(w, fmt.Sprintf("failed to decode the request: %s", err.Error()), http.StatusBadRequest)
		return
	}
	resp, err := handler.gameService.AnswerTheQuestion(r.Context(), gameId, questionId, req.Answer)
	if err != nil {
		handler.logger.Error("failed to answer the question", zap.Error(err))
		http.Error(w, fmt.Sprintf("failed to answer the question: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		handler.logger.Error("failed to encode the response", zap.Error(err))
		http.Error(w, fmt.Sprintf("failed to encode the response: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

// GetAllFlags — отладочный эндпоинт, возвращающий все флаги с их country_id.
// Используется фронтендовым роутом /debug/flags для визуальной проверки SVG.
func (handler *GameHandler) GetAllFlags(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

	flags := handler.gameService.GetAllFlags()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(flags); err != nil {
		handler.logger.Error("failed to encode flags response", zap.Error(err))
		http.Error(w, fmt.Sprintf("failed to encode response: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (handler *GameHandler) End(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

	gameIdS := mux.Vars(r)["gameId"]
	if gameIdS == "" {
		handler.logger.Error("Missing gameId in path")
		http.Error(w, "Missing gameId", http.StatusBadRequest)
		return
	}
	gameId, err := uuid.Parse(gameIdS)
	if err != nil {
		handler.logger.Error("Invalid gameId", zap.String("gameId", gameIdS))
		http.Error(w, "Invalid gameId", http.StatusBadRequest)
		return
	}

	questionsWithAnswer, err := handler.gameService.EndGame(r.Context(), gameId)
	if err != nil {
		handler.logger.Error("failed to end the game", zap.Error(err))
		http.Error(w, fmt.Sprintf("failed to end the game: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(questionsWithAnswer)
	if err != nil {
		handler.logger.Error("failed to encode the response", zap.Error(err))
		http.Error(w, fmt.Sprintf("failed to encode the response: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}
