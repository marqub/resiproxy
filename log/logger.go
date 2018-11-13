package log

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

//Logger wraps app and datatime for common log entries
func Logger() *log.Entry {
	return log.WithFields(log.Fields{"app": "resiproxy", "date": time.Now().UTC().Format(time.RFC3339)})
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})
}
