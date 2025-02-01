package log

import (
	"log"

	"github.com/fatih/color"
)

var (
	infolnLogger  *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
)

func Init(verbose bool) {
	if verbose {
		infolnLogger = log.New(color.Output, "", 0)
	}
	warningLogger = log.New(color.Error, color.YellowString("WARNING: "), 0)
	errorLogger = log.New(color.Error, color.RedString("ERROR: "), 0)
}

func Info(format string, v ...interface{}) {
	if infolnLogger != nil {
		infolnLogger.Printf(format, v...)
	}
}

func Warning(format string, v ...interface{}) {
	warningLogger.Printf(format, v...)
}

func Error(format string, v ...interface{}) {
	errorLogger.Printf(format, v...)
}

func Fatal(format string, v ...interface{}) {
	errorLogger.Fatalf(format, v...)
}
