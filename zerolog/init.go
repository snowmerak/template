package zerolog

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Kind string

const (
	KindEvent  Kind = "event"
	KindSystem Kind = "system"
)

type Type string

const (
	TypeSystemStart    Type = "system_start"
	TypeSystemShutdown Type = "system_shutdown"

	TypeEventUserLogin  Type = "user_login"
	TypeEventUserLogout Type = "user_logout"
)

type Level = zerolog.Level

const (
	LevelDebug Level = zerolog.DebugLevel
	LevelInfo  Level = zerolog.InfoLevel
	LevelWarn  Level = zerolog.WarnLevel
	LevelError Level = zerolog.ErrorLevel
	LevelFatal Level = zerolog.FatalLevel
	LevelPanic Level = zerolog.PanicLevel
)

var logger zerolog.Logger

func Init(appName string, level Level) {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.SetGlobalLevel(level)
	logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Str("app", appName).Logger()
}

func Logger() *zerolog.Logger {
	return &logger
}

func Debug() *zerolog.Event {
	return logger.Debug()
}

func Info() *zerolog.Event {
	return logger.Info()
}

func Warn() *zerolog.Event {
	return logger.Warn()
}

func Error() *zerolog.Event {
	return logger.Error()
}

func Send(event *zerolog.Event) {
	event.Send()
}

func Write(ev *zerolog.Event, kind Kind, typ Type) {
	ev.Str("kind", string(kind)).Str("type", string(typ)).Send()
}
