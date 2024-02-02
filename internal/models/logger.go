package models

import (
	"go.uber.org/zap"
)

func GetNewLogger() Logger {
	logger, _ := zap.NewProduction()
	sugarLogger := logger.Sugar()

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			sugarLogger.Errorf("Error during sync logger")
		}
	}(logger)

	return sugarLogger
}
