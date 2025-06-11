package logger

import (
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
