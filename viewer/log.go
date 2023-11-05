package viewer

import log "github.com/sirupsen/logrus"

func logInfo(s string, args ...interface{}) {
	log.Infof(s, args...)
}

func logDebug(s string, args ...interface{}) {
	log.Debugf(s, args...)
}

func logError(s string, args ...interface{}) {
	log.Errorf(s, args...)
}
