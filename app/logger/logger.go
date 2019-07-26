package logger

import (
	"go.uber.org/zap"
)

func New(mode string) (*zap.SugaredLogger, error) {
	var log *zap.Logger
	var err error

	if mode == "develop" {
		log, err = zap.NewDevelopment()
	} else {
		log, err = zap.NewProduction()
	}

	return log.Sugar(), err
}
