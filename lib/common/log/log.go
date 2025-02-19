package log

import (
	"fmt"
	"io"
	"log"

	"github.com/fatih/color"
)

var (
	infolnLogger       *log.Logger = nil
	warningLogger      *log.Logger = log.New(color.Error, color.YellowString("WARNING: "), 0)
	errorLogger        *log.Logger = log.New(color.Error, color.RedString("ERROR: "), 0)
	isQuiet            bool
	PanicInsteadOfExit bool
)

func Init(verbose bool, quiet bool) {
	if verbose {
		infolnLogger = log.New(color.Output, "INFO: ", 0)
	}
	isQuiet = quiet
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

func Printlnf(format string, v ...interface{}) {
	if !isQuiet {
		fmt.Printf(format, v...)
		fmt.Println()
	}
}

func Warning(format string, v ...interface{}) {
	warningLogger.Printf(format, v...)
}

func Error(format string, v ...interface{}) {
	errorLogger.Printf(format, v...)
}

func Fatal(format string, v ...interface{}) {
	if PanicInsteadOfExit {
		Error(format, v...)
		panic(fmt.Sprintf(format, v...))
	}
	errorLogger.Fatalf(format, v...)
}

func IsQuiet() bool {
	return isQuiet
}
