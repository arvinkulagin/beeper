package log

import (
	"log"
	"os"
)

type Logger struct {
	Out *log.Logger
	Err *log.Logger
}

func NewLogger() Logger {
	return Logger{
		Out: log.New(os.Stdout, "BEEPER ", log.LstdFlags),
		Err: log.New(os.Stderr, "ERROR ", log.LstdFlags|log.Lshortfile),
	}
}