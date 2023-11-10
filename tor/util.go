package tor

import (
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func GetIP() {
	c := http.DefaultClient
	res, err := c.Get("https://api.ipify.org")
	if err != nil {
		log.WithField("reason", err).Errorf("Could not fetch IP")
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.WithField("reason", err).Errorf("Could not fetch IP")
	}
	log.Infof("Current IP: %s", body)
}
