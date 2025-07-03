package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pythonistD/Guess-The-Flag/internal/schema"
	"github.com/pythonistD/Guess-The-Flag/internal/service/user"
	"go.uber.org/zap"
)

type UserHandler struct {
	userService user.Service
	logger      *zap.Logger
}

func NewUserHandler(userService user.Service, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

func (u *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Обрабатываем OPTIONS запросы для CORS
	if r.Method == "OPTIONS" {
		return
	}

	var req schema.Register
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		u.logger.Error("failed to unmarshal request body", zap.Error(err))
		http.Error(w, fmt.Sprintf("failed to unmarshal request body: %s", err.Error()), http.StatusBadRequest)
		return
	}
	token, err := u.userService.Register(r.Context(), req)
	if err != nil {
		u.logger.Error("failed to register user", zap.Error(err))
		http.Error(w, fmt.Sprintf("failed to register user: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(token)
	if err != nil {
		u.logger.Error("failed to marshal response", zap.Error(err))
		http.Error(w, fmt.Sprintf("failed to marshal response: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Обрабатываем OPTIONS запросы для CORS
	if r.Method == "OPTIONS" {
		return
	}

	var req schema.Login
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		u.logger.Error("failed to unmarshal request body", zap.Error(err))
		http.Error(w, fmt.Sprintf("failed to unmarshal request body: %s", err.Error()), http.StatusBadRequest)
		return
	}
	token, err := u.userService.Login(r.Context(), req)
	if err != nil {
		u.logger.Error("failed to login user", zap.Error(err))
		http.Error(w, fmt.Sprintf("failed to login user: %s", err.Error()), http.StatusForbidden)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(token)
	if err != nil {
		u.logger.Error("failed to marshal response", zap.Error(err))
		http.Error(w, fmt.Sprintf("failed to marshal response: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}
