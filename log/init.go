package log

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func Init(level uint32) {
	if level > 6 {
		fmt.Fprintf(os.Stderr, "Error: '--log' accepts numbers between 0 and 6, received %v\n", level)
		os.Exit(1)
	}

	var logLevel log.Level
	switch level {
	case 1:
		logLevel = log.FatalLevel
	case 2:
		logLevel = log.ErrorLevel
	case 3:
		logLevel = log.WarnLevel
	case 4:
		logLevel = log.InfoLevel
	case 5:
		logLevel = log.DebugLevel
	case 6:
		logLevel = log.TraceLevel
	}

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "02-01-2006 15:04:05 MST",
	})

	log.SetOutput(os.Stdout)
	log.SetLevel(logLevel)
}
