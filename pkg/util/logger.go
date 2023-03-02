package util

import (
	"sync"

	"go.uber.org/zap"
)

var logger *Logger
var once sync.Once

type Logger struct {
	*zap.SugaredLogger
}

func InitLogger() *Logger {
	once.Do(func() {
		l, _ := zap.NewDevelopment()
		logger = &Logger{
			l.Sugar(),
		}
	})

	return logger
}

func GetLogger() *Logger {
	return logger
}
