package log

import (
	"reflect"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

func TestLogger(t *testing.T) {
	tests := []struct {
		name string
		want *log.Entry
	}{
		{
			name: "ok",
			want: log.WithFields(log.Fields{"app": "resiproxy", "date": time.Now().UTC().Format(time.RFC3339)}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Logger(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Logger() = %v, want %v", got, tt.want)
			}
		})
	}
}
