package logger

import "log"

type Logger struct {
	Log bool
}

var logger *Logger

func New(log bool) {
	logger = &Logger{Log: log}
}

func Log(message string) {
	if logger.Log {
		log.Println(message)
	}
}
