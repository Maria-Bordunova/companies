package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Event struct {
	id      int
	message string
}

var (
	invalidArgMessage      = Event{1, "Invalid arg: %s"}
	invalidArgValueMessage = Event{2, "Invalid value for argument: %s: %v"}
	missingArgMessage      = Event{3, "Missing arg: %s"}
)

type Logger struct {
	*zap.SugaredLogger
}

type Config struct {
	Level string `env:"LOG_LEVEL" default:"debug"`
}

func New(cfg *Config) *Logger {
	zapLogger := newZap(cfg)

	logger := Logger{
		SugaredLogger: zapLogger.Sugar(),
	}
	return &logger
}

func newZap(cfg *Config) *zap.Logger {
	var logger *zap.Logger
	var err error

	logLevel := parseLevel(cfg.Level)

	zapCfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(logLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	logger, err = zapCfg.Build()

	if err != nil {
		panic(err)
	}

	return logger
}

func parseLevel(level string) zapcore.Level {
	logLevel := new(zapcore.Level)
	err := logLevel.Set(level)

	if err != nil {
		panic(err)
	}

	return *logLevel
}

func (l *Logger) InvalidArg(arg string) {
	l.SugaredLogger.Errorf(invalidArgMessage.message, arg)
}

func (l *Logger) InvalidArgValue(arg string, val string) {
	l.SugaredLogger.Errorf(invalidArgValueMessage.message, arg, val)
}

func (l *Logger) MissingArg(arg string) {
	l.SugaredLogger.Errorf(missingArgMessage.message, arg)
}

// With wraps the same function of SugaredLogger
func (l *Logger) With(args ...interface{}) *Logger {
	return &Logger{
		SugaredLogger: l.SugaredLogger.With(args...),
	}
}
