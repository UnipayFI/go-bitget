package log

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
)

type Logger interface {
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
	Debugf(format string, args ...any)
}

type defaultLogger struct {
	*slog.Logger
}

var (
	once sync.Once
	log  Logger = &defaultLogger{}
)

func GetDefaultLogger() Logger {
	once.Do(func() {
		errorOptions := &slog.HandlerOptions{
			Level: slog.LevelError,
		}
		log = &defaultLogger{
			Logger: slog.New(slog.NewTextHandler(os.Stdout, errorOptions)),
		}
	})
	return log
}

func (l *defaultLogger) Infof(format string, args ...any) {
	l.Info(fmt.Sprintf(format, args...))
}

func (l *defaultLogger) Warnf(format string, args ...any) {
	l.Warn(fmt.Sprintf(format, args...))
}

func (l *defaultLogger) Errorf(format string, args ...any) {
	l.Error(fmt.Sprintf(format, args...))
}

func (l *defaultLogger) Debugf(format string, args ...any) {
	l.Debug(fmt.Sprintf(format, args...))
}
