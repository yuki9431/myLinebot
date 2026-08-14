package logger

import (
	"bytes"
	"testing"
	"time"
)

func TestWriter(t *testing.T) {
	var buf bytes.Buffer
	logger := New(&buf)

	if logger == nil {
		t.Error("New is nil")
	} else {
		logger.Write("hello logger")

		dateLayout := "2006/01/02 15:04:05"
		n := time.Now()
		targetMessage := n.Format(dateLayout) + " hello logger\n"

		if buf.String() != targetMessage {
			t.Errorf("Error: '%s'", buf.String())
		}
	}
}
