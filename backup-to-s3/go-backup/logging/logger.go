package logging

import (
	"fmt"
	"log"
	"os"
)

type Logger interface {
	Error(err error)
	Warn(format string, args ...any)
	Info(format string, args ...any)
}

type logger struct {
	lgr *log.Logger
}

func New() Logger {
	log.SetPrefix("[backup-to-s3 addon] ")

	return &logger{
		lgr: log.New(os.Stdout, "[backup-to-s3 addon] ", log.Ldate|log.Ltime|log.Lmsgprefix),
	}
}

func (lg *logger) Error(err error) {
	lg.logLevelf("ERROR", "%v", err)
}

func (lg *logger) Warn(format string, args ...any) {
	lg.logLevelf("WARNING", format, args...)
}

func (lg *logger) Info(format string, args ...any) {
	lg.logLevelf("INFO", format, args...)
}

func (lg *logger) logLevelf(level, format string, args ...any) {
	lg.lgr.Printf("%s %s", level, fmt.Sprintf(format, args...))
}
