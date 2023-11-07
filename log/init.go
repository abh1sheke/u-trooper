package log

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func Init(level uint32) {
	if level > 6 || level < 1 {
		fmt.Fprintf(os.Stderr, "Error: '--log' accepts numbers between 0 and 6, received %v\n", level)
		os.Exit(1)
	}
	var logLevel = log.Level(level)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "02-01-2006 15:04:05 MST",
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(logLevel)
}
