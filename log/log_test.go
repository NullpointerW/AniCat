package log

import (
	"errors"
	"os"
	"testing"
	"time"
)

type TestSturct struct {
	Wall int
	Ext  int
}

func TestInfo(t *testing.T) {
	Infof(Struct{"my-test", 2, "hota l", "won z"}, "here is a testing %s", "no.1")
	Info(Struct{"my-test", 2, "hota l", "won z"}, "here is a testing")
	Info(Struct{"my-test", 2, "hota l", "won z"}, "here is a testing2")
	Info(Struct{"test", 2, "log.Info()"}, "here is a testing", " ", "2")
}

func TestError(t *testing.T) {
	err := errors.New("we got problem")

	Errorf(Struct{"error", err, "hota l", "won z"}, "")
	Error(Struct{"error", err, "hota l", "won z"}, "here is a testing")
	Error(Struct{"error", err, "hota l", "won z"}, "here is a testing2")
	Error(Struct{"error", err, "hota l", "won z"}, "here is a testing", " ", "2")
	logger := New("json", "mute", time.RFC3339Nano, os.Stderr)
	logger.Errorf(Struct{"error", err, "hota l", "won z"}, "")
	Error(Struct{"error", err, "hota l", "won z"}, "here is a testing")
	Error(Struct{"error", err, "hota l", "won z"}, "here is a testing2")
	Error(Struct{"error", err, "hota l", "won z"}, "here is a testing", " ", "2")
}
