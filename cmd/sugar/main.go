package main

import (
	"flag"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/yenole/sugar"
)

var (
	listen  *string = flag.String("listen", ":8080", "sugar listen")
	gListen *string = flag.String("glisten", ":8081", "sugar regist listen")
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
	sugar.New(logger).Run(&sugar.Option{
		Listen:  *listen,
		Gateway: *gListen,
	})
}
