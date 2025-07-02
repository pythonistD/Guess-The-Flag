package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/pythonistD/Guess-The-Flag/internal/schema"
	"github.com/pythonistD/Guess-The-Flag/internal/service/user"
	"go.uber.org/zap"
	"net/http"
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
	var req schema.Register
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		u.logger.Error("failed to unmarshal request body", zap.Error(err))
		http.Error(w, fmt.Sprintf("failed to unmarshal request body: %s", err.Error()), http.StatusBadRequest)
	}
	token, err := u.userService.Register(r.Context(), req)
	if err != nil {
		u.logger.Error("failed to register user", zap.Error(err))
		http.Error(w, fmt.Sprintf("failed to register user: %s", err.Error()), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(token)
	if err != nil {
		u.logger.Error("failed to marshal response", zap.Error(err))
		http.Error(w, fmt.Sprintf("failed to marshal response: %s", err.Error()), http.StatusInternalServerError)
	}
}

func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req schema.Login
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		u.logger.Error("failed to unmarshal request body", zap.Error(err))
		http.Error(w, fmt.Sprintf("failed to unmarshal request body: %s", err.Error()), http.StatusBadRequest)
	}
	token, err := u.userService.Login(r.Context(), req)
	if err != nil {
		u.logger.Error("failed to login user", zap.Error(err))
		http.Error(w, fmt.Sprintf("failed to login user: %s", err.Error()), http.StatusForbidden)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(token)
	if err != nil {
		u.logger.Error("failed to marshal response", zap.Error(err))
	}
}
