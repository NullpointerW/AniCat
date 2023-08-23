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
	logger := New("json", "mute", time.RFC3339Nano, true, os.Stderr)
	logger.Infof(nil, "this is a info msg")
}

func TestError(t *testing.T) {
	err := errors.New("we got problem")
	Errorf(Struct{"error", err, "hota l", "won z"}, "")
	Error(Struct{"error", err, "hota l", "won z"}, "here is a testing")
	Error(Struct{"error", err, "hota l", "won z"}, "here is a testing2")
	Error(Struct{"error", err, "hota l", "won z"}, "here is a testing", " ", "2")
	logger := New("json", "mute", time.RFC3339Nano, true, os.Stderr)
	logger.Errorf(Struct{"error", err, "hota l", "won z"}, "")
	Error(Struct{"error", err, "hota l", "won z"}, "here is a testing")
	Error(Struct{"error", err, "hota l", "won z"}, "here is a testing2")
	Error(Struct{"error", err, "hota l", "won z"}, "here is a testing", " ", "2")
}

func TestWarn(t *testing.T) {
	err := errors.New("we got problem")
	Warnf(Struct{"error", err, "hota l", "won z"}, "")
	Warn(Struct{"error", err, "hota l", "won z"}, "here is a testing")
	Warn(Struct{"error", err, "hota l", "won z"}, "here is a testing2")
	Warn(Struct{"error", err, "hota l", "won z"}, "here is a testing", " ", "2")
	logger := New("json", "warn", time.RFC3339Nano, true, os.Stderr)
	logger.Errorf(Struct{"error", err, "hota l", "won z"}, "")
	logger.Warnf(Struct{"error", err, "hota l", "won z"}, "")
	Error(Struct{"error", err, "hota l", "won z"}, "here is a testing")
	Error(Struct{"error", err, "hota l", "won z"}, "here is a testing2")
	Error(Struct{"error", err, "hota l", "won z"}, "here is a testing", " ", "2")
}

func TestDebug(t *testing.T) {
	Debug(Struct{"error", "<nil>", "hota l", "won z"}, "")
	Warn(Struct{"error", "<nil>", "hota l", "won z"}, "here is a testing")
	Warn(Struct{"error", "<nil>", "hota l", "won z"}, "here is a testing2")
	Warn(Struct{"error", "<nil>", "hota l", "won z"}, "here is a testing", " ", "2")
	Info(Struct{"error", "<nil>", "hota l", "won z"}, "you can not see me ", " ", "2")
	logger := New("json", "warn", time.RFC3339Nano, true, os.Stderr)
	logger.Errorf(Struct{"error", "<nil>", "hota l", "won z"}, "")
	logger.Warnf(Struct{"error", "<nil>", "hota l", "won z"}, "")
	Error(Struct{"error", "<nil>", "hota l", "won z"}, "here is a testing")
	Error(Struct{"error", "<nil>", "hota l", "won z"}, "here is a testing2")
	Error(Struct{"error", "<nil>", "hota l", "won z"}, "here is a testing", " ", "2")
}

func TestSlient(t *testing.T) {
	logger := New("json", "silent", time.RFC3339Nano, true, os.Stderr)
	logger.Errorf(Struct{"error", "<nil>", "hota l", "won z"}, "you can not see me")
	logger.Info(Struct{"error", "<nil>", "hota l", "won z"}, "you can not see me")
	logger.Debug(Struct{"error", "<nil>", "hota l", "won z"}, "you can not see me")
	logger.Warn(Struct{"error", "<nil>", "hota l", "won z"}, "you can not see me")
}

func TestSOURCE(t *testing.T) {
	logger := New("text", "info", time.RFC3339Nano, true, os.Stderr)
	logger.Errorf(Struct{"error", "<nil>", "hota l", "won z"}, "you can not see me")
	
}


