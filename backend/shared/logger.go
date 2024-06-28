package utils

import "github.com/sirupsen/logrus"

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	defer logger.Writer().Close()
}

func GetLogger() *logrus.Logger {
	return logger
}
