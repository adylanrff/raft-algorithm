package util

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func InitLogger(logPath string) {
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)

	f, err := os.OpenFile(logPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)
}
