package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

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
)

func main() {
	// Getting program args
	yamlConfigPath := flag.String("config", "../config.yml", "path to config file")
	logLevel := flag.String("log-level", "info", "logging level (debug, info, warn, error)")
	flag.Parse()

	// Устанавливаем уровень логирования через переменную окружения
	os.Setenv("LOG_LEVEL", *logLevel)

	logger := utils.NewLogger()
	r := mux.NewRouter()

	logger.Info("Starting Guess The Flag server")
	logger.Debug("Loading config", zap.String("config_path", *yamlConfigPath))

	cfg, err := config.LoadConfigFromFile(*yamlConfigPath)
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
		return
	}
	logger.Info("Config loaded successfully")

	database, err := db.NewPostgres(cfg.DBConfig)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
		return
	}
	logger.Info("Database connected successfully")

	// Init Storages
	countryStorage := storage.NewInMemoryCountryStorage(database)
	gameStorage := storage.NewInMemoryGameStorage()
	// Load countries into the storage
	err = countryStorage.InitCountryStorageState()
	if err != nil {
		logger.Error("Failed to load countries from database", zap.Error(err))
	} else {
		logger.Info("Countries loaded successfully")
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

	// Register Middleware and Routes - CORS должен быть первым!
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggerMiddleware(logger))
	r.Use(middleware.JWTMiddleware(jwtManager))

	authSubRouter := r.PathPrefix("/auth").Subrouter()
	authSubRouter.HandleFunc("/register", userHandler.Register).Methods("POST", "OPTIONS")
	authSubRouter.HandleFunc("/login", userHandler.Login).Methods("POST", "OPTIONS")

	gameSubRouter := r.PathPrefix("/game").Subrouter()
	gameSubRouter.HandleFunc("/start", gameHandler.Start).Methods("POST", "OPTIONS")
	gameSubRouter.HandleFunc("/{gameId}/questions/next", gameHandler.GetQuestion).Methods("POST", "OPTIONS")
	gameSubRouter.HandleFunc("/{gameId}/questions/{questionId}/answer", gameHandler.AnswerTheQuestion).Methods("POST", "OPTIONS")
	gameSubRouter.HandleFunc("/{gameId}/end", gameHandler.End).Methods("POST", "OPTIONS")

	// Публичный отладочный роут, отдающий все SVG флагов вместе с country_id.
	// Защита JWT отключена для него в middleware/auth.go.
	r.HandleFunc("/debug/flags", gameHandler.GetAllFlags).Methods("GET", "OPTIONS")

	logger.Info("Starting server",
		zap.String("address", cfg.Addr),
		zap.String("log_level", *logLevel),
	)

	server := &http.Server{
		Addr:         cfg.Addr,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Fatal("Server startup failed", zap.Error(err))
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	signal.Notify(c, os.Interrupt)
	logger.Info("Server started successfully. Press Ctrl+C to shutdown.")
	<-c

	logger.Info("Shutdown signal received")
	wait := time.Second * 5
	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	err = server.Shutdown(ctx)
	if err != nil {
		logger.Fatal("Server shutdown failed", zap.Error(err))
	}
	logger.Info("Server shut down gracefully")
	os.Exit(0)
}
