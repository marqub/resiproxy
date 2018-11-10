package log

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func Logger() *log.Entry {
	return log.WithFields(log.Fields{"appn": "resiproxy", "date": time.Now().UTC().Format(time.RFC3339)})
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})
}
