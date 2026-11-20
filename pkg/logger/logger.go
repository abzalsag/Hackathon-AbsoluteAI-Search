package logger

import (
	"log"
	"os"
)

type Logger struct {
	*log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

func (l *Logger) Info(msg string) {
	l.Printf("INFO: %s", msg)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.Printf("INFO: "+format, args...)
}

func (l *Logger) Error(msg string) {
	l.Printf("ERROR: %s", msg)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Printf("FATAL: "+format, args...)
	os.Exit(1)
}
