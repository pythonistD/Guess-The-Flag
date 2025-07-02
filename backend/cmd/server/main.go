package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pythonistD/Guess-The-Flag/internal/api/handlers"
	"github.com/pythonistD/Guess-The-Flag/internal/api/middleware"
	"github.com/pythonistD/Guess-The-Flag/internal/config"
	"github.com/pythonistD/Guess-The-Flag/internal/db"
	"github.com/pythonistD/Guess-The-Flag/internal/service/game"
	"github.com/pythonistD/Guess-The-Flag/internal/service/game/storage"
	"github.com/pythonistD/Guess-The-Flag/internal/service/user"
	"github.com/pythonistD/Guess-The-Flag/internal/utils"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// Getting program args
	yamlConfigPath := flag.String("config", "../../config.yaml", "path to config file")
	flag.Parse()

	logger := utils.NewLogger()
	r := mux.NewRouter()
	cfg, err := config.LoadConfigFromFile(*yamlConfigPath)
	logger.Debug("Loading config", zap.Any("config", cfg))
	if err != nil {
		logger.Fatal(err.Error())
		return
	}
	logger.Info("Config loaded successfully")
	database, err := db.NewPostgres(cfg.DBConfig)
	if err != nil {
		logger.Fatal(err.Error())
		return
	}
	logger.Info("Database loaded successfully")
	// Init Storages
	countryStorage := storage.NewInMemoryCountryStorage(database)
	gameStorage := storage.NewInMemoryGameStorage()
	// Load countries into the storage
	err = countryStorage.LoadCountriesFromDB()
	if err != nil {
		logger.Error(err.Error())
	}
	// Init security managers
	jwtManager := user.NewJWTManager(cfg.Secret, cfg.TokenLifetime)
	passManager := user.NewPasswordManager()
	// Init services
	userService := user.NewService(database, jwtManager, passManager)
	gameService := game.NewService(database, gameStorage, countryStorage)
	// Init handlers
	userHandler := handlers.NewUserHandler(userService, logger)
	gameHandler := handlers.NewGameHandler(gameService, logger)
	// Register Middleware and Routes
	r.Use(middleware.JWTMiddleware(jwtManager))
	r.Use(middleware.LoggerMiddleware(logger))

	authSubRouter := r.PathPrefix("/auth").Subrouter()
	authSubRouter.HandleFunc("/register", userHandler.Register).Methods("POST")
	authSubRouter.HandleFunc("/login", userHandler.Login).Methods("POST")

	gameSubRouter := r.PathPrefix("/game").Subrouter()
	gameSubRouter.HandleFunc("/start", gameHandler.Start).Methods("POST")
	gameSubRouter.HandleFunc("/question", gameHandler.GetQuestion).Methods("POST")
	gameSubRouter.HandleFunc("/answer", gameHandler.AnswerTheQuestion).Methods("POST")
	gameSubRouter.HandleFunc("/end", gameHandler.End).Methods("POST")
	logger.Info("Starting server at", zap.String("address", cfg.Addr))
	server := &http.Server{
		Addr:         cfg.Addr,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			err = fmt.Errorf("error starting server: %w", err)
			logger.Fatal(err.Error())
		}
	}()
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	signal.Notify(c, os.Interrupt)
	<-c
	wait := time.Second * 5
	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	err = server.Shutdown(ctx)
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Info("shutting down")
	os.Exit(0)
}
