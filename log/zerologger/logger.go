package zerologger

import (
	"os"

	"github.com/rezaAmiri123/edat/log"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

type zerologLogger struct {
	l zerolog.Logger
}

var _ log.Logger = (*zerologLogger)(nil)

func Logger(logger zerolog.Logger) log.Logger {
	zLog := logger.With().CallerWithSkipFrameCount(3).Logger()
	return &zerologLogger{l: zLog}
}
func NewZeroLogger(cfg log.Config) (log.Logger, error) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	var logger zerolog.Logger
	switch cfg.Environment {
	case "production":
		logger = zerolog.New(os.Stdout).
			Level(logLevelToZero(cfg.LogLevel)).
			With().
			Timestamp().
			Caller().
			Logger()
	default:
		stdout := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
			w.TimeFormat = "03:04:05.000PM"
		})
		logger = zerolog.New(stdout).
			Level(logLevelToZero(cfg.LogLevel)).
			With().
			Timestamp().
			Caller().
			Logger()
	}

	zLog := logger.With().CallerWithSkipFrameCount(3).Logger()
	return &zerologLogger{l: zLog}, nil
}

func (l *zerologLogger) Trace(msg string, fields ...log.Field) {
	if l.l.GetLevel() > zerolog.DebugLevel {
		return
	}
	logger := l.fields(l.l.With(), fields).Logger()
	logger.Debug().Msg(msg)
}

func (l *zerologLogger) Debug(msg string, fields ...log.Field) {
	if l.l.GetLevel() > zerolog.DebugLevel {
		return
	}
	logger := l.fields(l.l.With(), fields).Logger()
	logger.Debug().Msg(msg)

}

func (l *zerologLogger) Info(msg string, fields ...log.Field) {
	if l.l.GetLevel() > zerolog.InfoLevel {
		return
	}
	logger := l.fields(l.l.With(), fields).Logger()
	logger.Info().Msg(msg)
}

func (l *zerologLogger) Warn(msg string, fields ...log.Field) {
	if l.l.GetLevel() > zerolog.WarnLevel {
		return
	}
	logger := l.fields(l.l.With(), fields).Logger()
	logger.Warn().Msg(msg)

}

func (l *zerologLogger) Error(msg string, fields ...log.Field) {
	if l.l.GetLevel() > zerolog.ErrorLevel {
		return
	}
	logger := l.fields(l.l.With(), fields).Logger()
	logger.Error().Msg(msg)

}

func (l *zerologLogger) Sub(fields ...log.Field) log.Logger {
	return &zerologLogger{
		l: l.fields(l.l.With(), fields).Logger(),
	}
}

func (l *zerologLogger) fields(ctx zerolog.Context, fields []log.Field) zerolog.Context {
	for _, field := range fields {
		switch field.Type {
		case log.StringType:
			ctx = ctx.Str(field.Key, field.String)
		case log.IntType:
			ctx = ctx.Int(field.Key, field.Int)
		case log.DurationType:
			ctx = ctx.Str(field.Key, field.Duration.String())
		case log.ErrorType:
			ctx = ctx.Stack().Err(field.Error)
		}
	}

	return ctx
}

func logLevelToZero(level log.Level) zerolog.Level {
	switch level {
	case log.PANIC:
		return zerolog.PanicLevel
	case log.ERROR:
		return zerolog.ErrorLevel
	case log.WARN:
		return zerolog.WarnLevel
	case log.INFO:
		return zerolog.InfoLevel
	case log.DEBUG, log.TRACE:
		return zerolog.DebugLevel
	default:
		return zerolog.InfoLevel

	}
}
