package log

import (
	// "fmt"
	"errors"
	"os"
	"testing"
	// "log/slog"
)

type TestSturct struct {
	Wall int
	Ext  int
}

func TestInfo(t *testing.T) {
	Init("text", "error", "2006-01-02T15:04:05.999", os.Stderr)
	Infof(Struct{"my-test", 2, "hota l", "won z"}, "here is a testing %s", "no.1")
	Info(Struct{"my-test", 2, "hota l", "won z"}, "here is a testing")
	Info(Struct{"my-test", 2, "hota l", "won z"}, "here is a testing2")
	Info(Struct{"test", 2, "log.Info()"}, "here is a testing", " ", "2")
}

func TestError(t *testing.T) {
	err := errors.New("we got problem")
	Init("text", "error", "2006-01-02T15:04:05.999", os.Stderr)
	Errorf(Struct{"error", err, "hota l", "won z"}, "")
	Error(Struct{"error", err, "hota l", "won z"}, "here is a testing")
	Error(Struct{"error", err, "hota l", "won z"}, "here is a testing2")
	Error(Struct{"error", err, "hota l", "won z"}, "here is a testing", " ", "2")
}
