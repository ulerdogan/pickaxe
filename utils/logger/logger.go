package logger_utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
)

var (
	log LoggerUtil
)

type LoggerUtil struct {
	logger zerolog.Logger
}

func init() {
	log.logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC822,
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("||> %s", i)
		},
		FormatCaller: func(i interface{}) string {
			return filepath.Base(fmt.Sprintf("%s", i))
		},
		FormatErrFieldName: func(i interface{}) string {
			return fmt.Sprintf("\n%s: ", i)
		},
	}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Logger()
}

func Info(msg string) {
	log.logger.Info().Msg(msg)
}

func Trace(msg string) {
	log.logger.Trace().Msg(msg)
}

func Error(err error, msg string) {
	log.logger.Error().Err(err).Msg(msg)
}
