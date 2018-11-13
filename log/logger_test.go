package log

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"reflect"
	"testing"
	"time"
)

var testLogger *logrus.Entry
var testWriter *bytes.Buffer

func setupTestCase(t *testing.T) func(t *testing.T) {
	testWriter = bytes.NewBuffer([]byte{})
	logrus.SetOutput(testWriter)
	return func(t *testing.T) {
	}
}

func TestLogger(t *testing.T) {
	teardown := setupTestCase(t)
	defer teardown(t)
	now := time.Now()
	tests := []struct {
		name string
		want []byte
	}{
		{
			name: "ok",
			want: []byte(fmt.Sprintf(`{"app":"resiproxy","date":"%v","level":"error","msg":"Test Error","time":"%v"}
`, now.UTC().Format(time.RFC3339), now.Format(time.RFC3339))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Logger(); got != nil {
				// Test out the logger to make sure it works right
				Logger().WithTime(time.Now()).Error("Test Error")
				if data := testWriter.Bytes(); !reflect.DeepEqual(data, tt.want) {
					t.Errorf("Logger() didn't write the correct output = %v, want %v", string(data), string(tt.want))
				}

			} else {
				t.Error("No Logger returned")
			}
		})
	}
}
