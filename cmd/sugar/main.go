package main

import (
	"flag"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/yenole/sugar"
)

func main() {
	flag.Parse()

	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.Formatter = &logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	}
	logger.SetLevel(logrus.DebugLevel)
	sugar.New(logger).Run()
}
