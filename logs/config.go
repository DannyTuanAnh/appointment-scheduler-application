package logs

import (
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

func InitLogger(path string) zerolog.Logger {
	logger := zerolog.New(&lumberjack.Logger{
		Filename:   path,
		MaxSize:    1, // megabytes
		MaxBackups: 5,
		MaxAge:     5,    //days
		Compress:   true, // disable by default
		LocalTime:  true, // use local time for timestamps
	}).With().Timestamp().Logger()

	return logger
}
