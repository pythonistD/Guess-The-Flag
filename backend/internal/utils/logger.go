package utils

import (
	"go.uber.org/zap"
	"log"
	"os"
)

var (
	LogInfo  *log.Logger
	LogWarn  *log.Logger
	LogError *log.Logger
)

func Init() {
	LogInfo = log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lshortfile)
	LogWarn = log.New(os.Stdout, "[WARNING] ", log.LstdFlags|log.Lshortfile)
	LogError = log.New(os.Stdout, "[ERROR] ", log.LstdFlags|log.Lshortfile)
}

type Config struct {
	Level       string
	Environment string
	OutputPaths []string
}

func getOrDefault(env string, defaultValue string) string {
	if val := os.Getenv(env); val != "" {
		return val
	}
	return defaultValue
}

func getLoggerConfig() Config {
	return Config{
		Level:       getOrDefault("LOG_LEVEL", "info"),
		Environment: getOrDefault("ENVIRONMENT", "development"),
		OutputPaths: []string{"stdout"},
	}
}

func NewLogger() *zap.Logger {
	config := getLoggerConfig()

	var zapConfig zap.Config
	if config.Environment == "production" {
		zapConfig = zap.NewProductionConfig()
	} else {
		zapConfig = zap.NewDevelopmentConfig()
	}

	switch config.Level {
	case "debug":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	if len(config.OutputPaths) > 0 {
		zapConfig.OutputPaths = config.OutputPaths
	}
	logger, err := zapConfig.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.ErrorLevel),
	)
	if err != nil {
		fallbackLogger, _ := zap.NewProduction()
		fallbackLogger.Error("Failed to initialize zap utils", zap.Error(err))
		return fallbackLogger
	}
	return logger
}
