// Package unit provides ...
package unit

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/yenole/sugar/logger"
)

func newLogger() logger.Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.Formatter = &logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	}
	logger.SetLevel(logrus.DebugLevel)
	return logger
}
