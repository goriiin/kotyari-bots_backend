package logger

import (
	"os"

	"github.com/goriiin/kotyari-bots_backend/pkg/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

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

func (l *Logger) Debug() *zerolog.Event { return l.log.Debug() }
func (l *Logger) Info() *zerolog.Event  { return l.log.Info() }
func (l *Logger) Warn() *zerolog.Event  { return l.log.Warn() }
func (l *Logger) Error() *zerolog.Event { return l.log.Error() }
func (l *Logger) Fatal() *zerolog.Event { return l.log.Fatal() }
