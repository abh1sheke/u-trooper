package log

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func Init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		TimestampFormat:        "02-01-2006 15:04:05 MST",
	})
	log.SetOutput(os.Stdout)
  log.SetLevel(log.DebugLevel)
}
