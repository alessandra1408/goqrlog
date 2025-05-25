package log

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Log interface {
	Level() zapcore.Level
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Fatal(args ...interface{})
	With(args ...interface{}) *zap.SugaredLogger
	Sync() error
}

func New(z *zap.SugaredLogger) Log {
	return &log{
		z: z,
	}
}

type log struct {
	z *zap.SugaredLogger
}

func (l *log) Level() zapcore.Level {
	return l.z.Level()
}

func (l *log) With(args ...interface{}) *zap.SugaredLogger {
	return l.z.With(args...)
}

func (l *log) Info(args ...interface{}) {
	l.z.Info(args...)
}

func (l *log) Infof(template string, args ...interface{}) {
	l.z.Infof(template, args...)
}

func (l *log) Warn(args ...interface{}) {
	l.z.Warn(args...)
}

func (l *log) Warnf(template string, args ...interface{}) {
	l.z.Warnf(template, args...)
}

func (l *log) Error(args ...interface{}) {
	l.z.Error(args...)
}

func (l *log) Errorf(template string, args ...interface{}) {
	l.z.Errorf(template, args...)
}

func (l *log) Fatal(args ...interface{}) {
	l.z.Fatal(args...)
}

func (l *log) Sync() error {
	return l.z.Sync()
}

func LogWithUserAgent(log Log, userAgent string) Log {
	if userAgent == "" {
		return log.With("user-agent", "empty")
	}

	for _, e := range strings.Split(userAgent, ";") {
		if e != "" && strings.Contains(e, "=") {
			s := strings.Split(e, "=")
			if len(s) == 2 {
				log = log.With("ua-"+s[0], s[1])
			}
		}
	}

	return log.With("user-agent", userAgent)
}
