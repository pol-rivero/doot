package log

import (
	"io"
	"log"

	"github.com/fatih/color"
)

var (
	infolnLogger  *log.Logger = nil
	warningLogger *log.Logger = log.New(color.Error, color.YellowString("WARNING: "), 0)
	errorLogger   *log.Logger = log.New(color.Error, color.RedString("ERROR: "), 0)
)

func Init(verbose bool, quiet bool) {
	if verbose {
		infolnLogger = log.New(color.Output, "", 0)
	}
	if quiet {
		warningLogger.SetOutput(io.Discard)
		errorLogger.SetOutput(io.Discard)
	}
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
