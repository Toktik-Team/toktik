package logging

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})

	log.SetOutput(os.Stdout)

	if os.Getenv("ENV") == "prod" {
		log.SetLevel(log.WarnLevel)
	} else {
		log.SetLevel(log.DebugLevel)
	}
}

var env = os.Getenv("ENV")

// Logger Add fields you want to log by default.
var Logger = log.WithFields(log.Fields{
	"ENV": env,
})
