package log

import (
	"io"
	// "os"
	eslog "github.com/NullpointerW/anicat/pkg/log"
	"sync"
)

var (
	logger *eslog.EnhanceLogger
	_init  sync.Once
)

type Struct []any

func (s *Struct) Append(a ...any) {
	*s = append(*s, a...)
}

func (s *Struct) Clear() {
	*s = nil
}

func NewUrlStruct(a ...any) Struct {
	s := Struct{"url"}
	return append(s, a...)
}

// Init initialization a structure log  must use Init if imported, otherwise trigger nil pointer panic
// Init initializes the logging system with the specified parameters.
// It ensures that the initialization is performed only once.
//
// Parameters:
//   - handleType: The type of log handling (e.g., console, file).
//   - level: The logging level (e.g., debug, info, warn, error).
//   - timeLayout: The time format layout for log timestamps.
//   - sourceFile: A boolean indicating whether to include the source file information in the log.
//   - out: An io.Writer where the logs will be written.
func Init(handleType, level, timeLayout string, sourceFile bool, out io.Writer) {
	_init.Do(func() {
		logger = eslog.New(handleType, level, timeLayout, sourceFile, 1, out)
	})
}

func Info(s Struct, a ...any) {
	logger.Info(eslog.Struct(s), eslog.Msg(a...))
}

func Infof(s Struct, format string, a ...any) {
	logger.Info(eslog.Struct(s), eslog.Msgf(format, a...))
}

func Warn(s Struct, a ...any) {
	logger.Warn(eslog.Struct(s), eslog.Msg(a...))
}

func Warnf(s Struct, format string, a ...any) {
	logger.Warn(eslog.Struct(s), eslog.Msgf(format, a...))
}

func Debug(s Struct, a ...any) {
	logger.Debug(eslog.Struct(s), eslog.Msg(a...))
}

func Debugf(s Struct, format string, a ...any) {
	logger.Debug(eslog.Struct(s), eslog.Msgf(format, a...))
}

func Error(s Struct, a ...any) {
	logger.Error(eslog.Struct(s), eslog.Msg(a...))
}
func Errorf(s Struct, format string, a ...any) {
	logger.Error(eslog.Struct(s), eslog.Msgf(format, a...))
}
