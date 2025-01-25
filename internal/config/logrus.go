package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	log := logrus.New()

	level, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		level = logrus.DebugLevel
	}

	log.SetLevel(logrus.Level(level))
	log.SetFormatter(&logrus.JSONFormatter{})

	return log
}
