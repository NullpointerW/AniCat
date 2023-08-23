package log

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"
)

var defaultLogger *EnhanceLogger

func init() {
	defaultLogger = New("text", "Info", time.RFC3339, false, os.Stderr)
}

type EnhanceLogger struct {
	inner *slog.Logger
}

type Struct []any

func slogLevel(l string) slog.Level {
	switch strings.ToLower(l) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "err", "error":
		return slog.LevelError
	case "silent", "mute":
		return slog.LevelError + 1
	default:
		return slog.LevelInfo
	}
}

func New(handleType, level, timeLayout string, shortSourceFile bool, out io.Writer) *EnhanceLogger {
	timeAttrFunc := timeFormat(timeLayout)
	opt := &slog.HandlerOptions{
		AddSource:   shortSourceFile,
		ReplaceAttr: timeAttrFunc,
		Level:       slogLevel(level),
	}
	var handler slog.Handler
	switch strings.ToLower(handleType) {
	case "json":
		handler = slog.NewJSONHandler(out, opt)
	default:
		handler = slog.NewTextHandler(out, opt)
	}
	return &EnhanceLogger{
		slog.New(handler),
	}
}

func timeFormat(layout string) func(groups []string, a slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			t := a.Value.Time()
			a.Value = slog.StringValue(t.Format(layout))
		}
		if a.Key == slog.SourceKey {
			if src, ok := a.Value.Any().(*slog.Source); ok {
				shortPath := ""
				fullPath := src.File
				seps := strings.Split(fullPath, "/")
				shortPath += seps[len(seps)-1]
				shortPath += fmt.Sprintf(":%d", src.Line)
				a.Value = slog.StringValue(shortPath)
			}
		}
		return a
	}
}

func Msg(a ...any) string {
	return fmt.Sprint(a...)
}

func Msgf(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}

func (el *EnhanceLogger) Debug(s Struct, a ...any) {
	el.inner.Debug(Msg(a...), s...)
}

func (el *EnhanceLogger) Debugf(s Struct, format string, a ...any) {
	el.inner.Debug(Msgf(format, a...), s...)
}

func (el *EnhanceLogger) Info(s Struct, a ...any) {
	el.inner.Info(Msg(a...), s...)
}

func (el *EnhanceLogger) Infof(s Struct, format string, a ...any) {
	el.inner.Info(Msgf(format, a...), s...)
}

func (el *EnhanceLogger) Warn(s Struct, a ...any) {
	el.inner.Warn(Msg(a...), s...)
}

func (el *EnhanceLogger) Warnf(s Struct, format string, a ...any) {
	el.inner.Warn(Msgf(format, a...), s...)
}

func (el *EnhanceLogger) Error(s Struct, a ...any) {
	el.inner.Error(Msg(a...), s...)
}

func (el *EnhanceLogger) Errorf(s Struct, format string, a ...any) {
	el.inner.Error(Msgf(format, a...), s...)
}

func Debug(s Struct, a ...any) {
	defaultLogger.inner.Debug(Msg(a...), s...)
}

func Debugf(s Struct, format string, a ...any) {
	defaultLogger.inner.Debug(Msgf(format, a...), s...)
}

// Spaces will not like the fmt.Println() that added between operands and a newline is appended.
// add spaces between operands to achieve the same effect:
//
//	Info(Struct{"test","log.Info()"},"here is a testing", " ", "log.Info()")
func Info(s Struct, a ...any) {
	defaultLogger.inner.Info(Msg(a...), s...)
}

func Infof(s Struct, format string, a ...any) {
	defaultLogger.inner.Info(Msgf(format, a...), s...)
}

func Warn(s Struct, a ...any) {
	defaultLogger.inner.Warn(Msg(a...), s...)
}

func Warnf(s Struct, format string, a ...any) {
	defaultLogger.inner.Warn(Msgf(format, a...), s...)
}

func Error(s Struct, a ...any) {
	defaultLogger.inner.Error(Msg(a...), s...)
}

func Errorf(s Struct, format string, a ...any) {
	defaultLogger.inner.Error(Msgf(format, a...), s...)
}
