package config

import (
	"github.com/sirupsen/logrus"
	"os"
)

func InitLogger() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}
