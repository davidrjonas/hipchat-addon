package addon

import (
	"log"
	"os"
)

type AddonLogger interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}

type StdLogger struct {
	Logger *log.Logger
}

func NewStdLogger() *StdLogger {
	return &StdLogger{
		Logger: log.New(os.Stderr, "hipchat: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (log *StdLogger) Info(args ...interface{}) {
	log.Logger.Print(args...)
}

func (log *StdLogger) Infof(format string, args ...interface{}) {
	log.Logger.Printf(format, args...)
}

func (log *StdLogger) Error(args ...interface{}) {
	log.Logger.Print(args...)
}

func (log *StdLogger) Errorf(format string, args ...interface{}) {
	log.Logger.Printf("Error: "+format, args...)
}
