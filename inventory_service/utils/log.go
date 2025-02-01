package utils

import (
	"go.uber.org/zap"
)

func InitLogger() (*zap.Logger, error) {

	Logger, err := zap.NewDevelopment()

	if err != nil {
		return nil, err
	}

	return Logger, nil
}
