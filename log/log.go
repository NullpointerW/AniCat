package log

import (
	"fmt"
	"io"
	"log/slog"

	// "os"
	"strings"
	"sync"
)

var (
	logger *slog.Logger
	init_  sync.Once
)

func slogLevel(l string) slog.Level {
	switch strings.ToLower(l) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "err", "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

type Struct []any

func Init(handleType, level, timeLayout string, out io.Writer) {
	init_.Do(func() {
		logInit(handleType, level, timeLayout, out)
	})
}

func logInit(handleType, level, timeLayout string, out io.Writer) {
	timeAttrFunc := timeFormat(timeLayout)
	opt := &slog.HandlerOptions{
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
	logger = slog.New(handler)
}

func timeFormat(layout string) func(groups []string, a slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			t := a.Value.Time()
			a.Value = slog.StringValue(t.Format(layout))
		}
		if a.Key == slog.MessageKey {
			a.Key = ""
		}
		return a
	}
}

// func init() {
// 	timeAttrFunc := timeFormat("2006-01-02T15:04:05.999")
// 	opt := &slog.HandlerOptions{
// 		ReplaceAttr: timeAttrFunc,
// 		Level:       level,
// 	}
// 	logger = slog.New(slog.NewTextHandler(os.Stderr, opt))
// }

func Msg(a ...any) string {
	String := fmt.Sprint(a...)
	return strings.TrimSuffix(String, "\n")
}

func Msgf(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}

// Spaces will not like the fmt.Println() that added between operands and a newline is appended.
// add spaces between operands to achieve the same effect:
//
//	Info(Struct{"test","log.Info()"},"here is a testing", " ", "log.Info()")
func Info(s Struct, a ...any) {
	logger.Info(Msg(a...), s...)
}

func Infof(s Struct, format string, a ...any) {
	logger.Info(Msgf(format, a...), s...)
}

func Warn(s Struct, a ...any) {
	logger.Warn(Msg(a...), s...)
}

func Warnf(s Struct, format string, a ...any) {
	logger.Warn(Msgf(format, a...), s...)
}

func Debug(s Struct, a ...any) {
	logger.Debug(Msg(a...), s...)
}

func Debugf(s Struct, format string, a ...any) {
	logger.Debug(Msgf(format, a...), s...)
}

func Error(s Struct, a ...any) {
	logger.Error(Msg(a...), s...)
}

func Errorf(s Struct, format string, a ...any) {
	logger.Error(Msgf(format, a...), s...)
}
