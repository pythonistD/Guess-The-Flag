package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pythonistD/Guess-The-Flag/internal/config"
	"github.com/pythonistD/Guess-The-Flag/internal/db"
	"github.com/pythonistD/Guess-The-Flag/internal/db/repo"
	"github.com/pythonistD/Guess-The-Flag/internal/logger"
	"net/http"
	"time"
)

func main() {
	logger.Init()
	r := mux.NewRouter()
	cfg, err := config.LoadConfigFromFile("./config.yml")
	if err != nil {
		logger.LogError.Fatal(err)
		return
	}
	logger.LogInfo.Println("Config loaded successfully")
	database, err := db.NewPostgres(cfg.DBConfig)
	if err != nil {
		logger.LogError.Fatal(err)
		return
	}
	logger.LogInfo.Println("Database loaded successfully")
	usersRepo := repo.NewUsersRepo(database)

	logger.LogInfo.Printf("Starting server at %s\n", cfg.Addr)
	server := &http.Server{
		Addr:         cfg.Addr,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err = server.ListenAndServe()
	if err != nil {
		err = fmt.Errorf("error starting server: %w", err)
		logger.LogError.Fatal(err)
		return
	}
}
