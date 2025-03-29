package zerolog

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Kind string

const (
	KindEvent   Kind = "event"
	KindSystem  Kind = "system"
	KindMeasure Kind = "measure"
)

type Status string

const (
	StatusSuccess    Status = "success"
	StatusFailure    Status = "failure"
	StatusCancelled  Status = "cancelled"
	StatusInProgress Status = "in_progress"
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

func New(status Status, level Level) *zerolog.Event {
	ev := (*zerolog.Event)(nil)
	switch level {
	case LevelDebug:
		ev = logger.Debug()
	case LevelInfo:
		ev = logger.Info()
	case LevelWarn:
		ev = logger.Warn()
	case LevelError:
		ev = logger.Error()
	case LevelFatal:
		ev = logger.Fatal()
	case LevelPanic:
		ev = logger.Panic()
	default:
		ev = logger.Info()
	}

	switch status {
	case StatusSuccess:
		return ev.Str("status", string(StatusSuccess))
	case StatusFailure:
		return ev.Str("status", string(StatusFailure))
	case StatusCancelled:
		return ev.Str("status", string(StatusCancelled))
	case StatusInProgress:
		return ev.Str("status", string(StatusInProgress))
	default:
		return ev.Str("status", string(StatusSuccess))
	}
}

func Send(event *zerolog.Event) {
	event.Send()
}

func Write(ev *zerolog.Event, kind Kind, typ Type) {
	ev.Str("kind", string(kind)).Str("type", string(typ)).Send()
}
