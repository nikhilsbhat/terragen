package cmd

import (
	goCdLogger "github.com/nikhilsbhat/gocd-sdk-go/pkg/logger"
	"github.com/sirupsen/logrus"
)

var (
	cliLogLevel string
	cliLogger   *logrus.Logger
)

func InitLogger(logLevel string) {
	logger := logrus.New()
	logger.SetLevel(goCdLogger.GetLoglevel(logLevel))
	logger.WithField("terragen", true)
	logger.SetFormatter(&logrus.JSONFormatter{})
	cliLogger = logger
}
