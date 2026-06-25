package logger

import (
	stderrors "errors"
	"os"
	"strings"

	"github.com/goriiin/kotyari-bots_backend/pkg/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

// joinErrs combines any number of (possibly nil) errors into a single error,
// returning nil when none are non-nil. Used so the variadic logging helpers
// below attach every supplied error instead of silently dropping all of them
// whenever more than one is passed.
func joinErrs(errs ...error) error {
	return stderrors.Join(errs...)
}

const service = "service_name"

type Logger struct {
	log zerolog.Logger
}

func NewLogger(serviceName string, appCfg *config.ConfigBase) *Logger {
	var level zerolog.Level
	if appCfg.IsProduction() {
		level = zerolog.InfoLevel
	} else {
		level = zerolog.DebugLevel
	}

	zerolog.SetGlobalLevel(level)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str(service, serviceName).
		Caller().
		Logger()

	return &Logger{logger}
}

func (l *Logger) logWithDetails(event *zerolog.Event, err error, withStack bool, msg ...string) {
	if withStack {
		event = event.Stack()
	}
	if err != nil {
		event = event.Err(err)
	}

	if len(msg) > 0 {
		event.Msg(strings.Join(msg, " "))
	} else {
		event.Send()
	}
}

func (l *Logger) Fatal(err error, withStack bool, msg ...string) {
	l.logWithDetails(l.log.Fatal(), err, withStack, msg...)
}

func (l *Logger) Error(err error, withStack bool, msg ...string) {
	l.logWithDetails(l.log.Error(), err, withStack, msg...)
}

func (l *Logger) Warn(msg string, err ...error) {
	if joined := joinErrs(err...); joined != nil {
		l.log.Warn().Err(joined).Msg(msg)
		return
	}

	l.log.Warn().Msg(msg)
}

func (l *Logger) Info(msg string, err ...error) {
	if joined := joinErrs(err...); joined != nil {
		l.log.Info().Err(joined).Msg(msg)
		return
	}

	l.log.Info().Msg(msg)
}

func (l *Logger) Debug(msg string, err ...error) {
	if joined := joinErrs(err...); joined != nil {
		l.log.Debug().Err(joined).Msg(msg)
		return
	}

	l.log.Debug().Msg(msg)
}
