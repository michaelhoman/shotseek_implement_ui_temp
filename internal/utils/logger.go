package utils

import "go.uber.org/zap"

// Logger is a global logger instance
var Logger *zap.SugaredLogger

// InitLogger initializes the logger and assigns it to the global variable
func InitLogger() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	Logger = logger
}

// CleanupLogger flushes logs before the program exits
func CleanupLogger() {
	if Logger != nil {
		Logger.Sync()
	}
}
