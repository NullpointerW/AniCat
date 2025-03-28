package log_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/NullpointerW/anicat/log"
)

func TestDebug(t *testing.T) {
	log.Init("text", "debug", "2006年01月02日", true, os.Stderr)
	log.Debug(log.Struct{"test", "debug"}, "hello", "world")
	log.Debug(log.Struct{"test", "debug"}, "hello", " ", "world")
	log.Debugf(log.Struct{"test", "debug"}, "hello %s", " world")
	log.Info(log.Struct{"test", "info"}, "hello", "world")
	log.Info(log.Struct{"test", "info"}, "hello", " ", "world")
	log.Infof(log.Struct{"test", "info"}, "hello %s", " world")
	log.Warn(log.Struct{"test", "warn"}, "hello", "world")
	log.Warn(log.Struct{"test", "warn"}, "hello", " ", "world")
	log.Warnf(log.Struct{"test", "warn"}, "hello %s", " world")
	log.Error(log.Struct{"test", "error"}, "hello", "world")
	log.Error(log.Struct{"test", "error"}, "hello", " ", "world")
	log.Errorf(log.Struct{"test", "error"}, "hello %s", " world")
}

func TestInfo(t *testing.T) {
	log.Init("text", "info", "2006年01月02日", true, os.Stderr)
	log.Debug(log.Struct{"test", "debug"}, "hello", "world")
	log.Debug(log.Struct{"test", "debug"}, "hello", " ", "world")
	log.Debugf(log.Struct{"test", "debug"}, "hello %s", " world")
	log.Info(log.Struct{"test", "info"}, "hello", "world")
	log.Info(log.Struct{"test", "info"}, "hello", " ", "world")
	log.Infof(log.Struct{"test", "info"}, "hello %s", " world")
	log.Warn(log.Struct{"test", "warn"}, "hello", "world")
	log.Warn(log.Struct{"test", "warn"}, "hello", " ", "world")
	log.Warnf(log.Struct{"test", "warn"}, "hello %s", " world")
	log.Error(log.Struct{"test", "error"}, "hello", "world")
	log.Error(log.Struct{"test", "error"}, "hello", " ", "world")
	log.Errorf(log.Struct{"test", "error"}, "hello %s", " world")
}
func TestWarn(t *testing.T) {
	log.Init("text", "Warn", "2006年01月02日", true, os.Stderr)
	log.Debug(log.Struct{"test", "debug"}, "hello", "world")
	log.Debug(log.Struct{"test", "debug"}, "hello", " ", "world")
	log.Debugf(log.Struct{"test", "debug"}, "hello %s", " world")
	log.Info(log.Struct{"test", "info"}, "hello", "world")
	log.Info(log.Struct{"test", "info"}, "hello", " ", "world")
	log.Infof(log.Struct{"test", "info"}, "hello %s", " world")
	log.Warn(log.Struct{"test", "warn"}, "hello", "world")
	log.Warn(log.Struct{"test", "warn"}, "hello", " ", "world")
	log.Warnf(log.Struct{"test", "warn"}, "hello %s", " world")
	log.Error(log.Struct{"test", "error"}, "hello", "world")
	log.Error(log.Struct{"test", "error"}, "hello", " ", "world")
	log.Errorf(log.Struct{"test", "error"}, "hello %s", " world")
}
func TestError(t *testing.T) {
	log.Init("text", "err", "2006年01月02日", true, os.Stderr)
	log.Debug(log.Struct{"test", "debug"}, "hello", "world")
	log.Debug(log.Struct{"test", "debug"}, "hello", " ", "world")
	log.Debugf(log.Struct{"test", "debug"}, "hello %s", " world")
	log.Info(log.Struct{"test", "info"}, "hello", "world")
	log.Info(log.Struct{"test", "info"}, "hello", " ", "world")
	log.Infof(log.Struct{"test", "info"}, "hello %s", " world")
	log.Warn(log.Struct{"test", "warn"}, "hello", "world")
	log.Warn(log.Struct{"test", "warn"}, "hello", " ", "world")
	log.Warnf(log.Struct{"test", "warn"}, "hello %s", " world")
	log.Error(log.Struct{"test", "error"}, "hello", "world")
	log.Error(log.Struct{"test", "error"}, "hello", " ", "world")
	log.Errorf(log.Struct{"test", "error"}, "hello %s", " world")
}
func TestMute(t *testing.T) {
	log.Init("text", "mute", "2006年01月02日", true, os.Stderr)
	log.Debug(log.Struct{"test", "debug"}, "hello", "world")
	log.Debug(log.Struct{"test", "debug"}, "hello", " ", "world")
	log.Debugf(log.Struct{"test", "debug"}, "hello %s", " world")
	log.Info(log.Struct{"test", "info"}, "hello", "world")
	log.Info(log.Struct{"test", "info"}, "hello", " ", "world")
	log.Infof(log.Struct{"test", "info"}, "hello %s", " world")
	log.Warn(log.Struct{"test", "warn"}, "hello", "world")
	log.Warn(log.Struct{"test", "warn"}, "hello", " ", "world")
	log.Warnf(log.Struct{"test", "warn"}, "hello %s", " world")
	log.Error(log.Struct{"test", "error"}, "hello", "world")
	log.Error(log.Struct{"test", "error"}, "hello", " ", "world")
	log.Errorf(log.Struct{"test", "error"}, "hello %s", " world")
}
func 	TestXxx(t *testing.T) {
	ax :=0

	type s struct {a int}
	var b s
	b.a=1
	type ss s
	var bb ss
	bb.a=1
	bb=ss(b)
	b=s(bb)
	defer func ()  {
		fmt.Println(ax)
	}()
	ax=1
}