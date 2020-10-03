package utils

import (
	"github.com/google/logger"
	"os"
)

const verbose = false

var loggerFile *os.File

func LoggerSetup() {
	var err error
	loggerFile, err = os.OpenFile(LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
	}
	logger.Init("LoggerExample", verbose, false, loggerFile)
}

func LoggerClose() {
	loggerFile.Close()
	logger.Close()
}
