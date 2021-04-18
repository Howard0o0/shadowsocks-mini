package tinylog

import (
	"fmt"
	"log"
	"os"
)

// logger level
const (
	DebugLevel int = iota
	InfoLevel
	ErrorLevel
	FatalLevel
)

type loggerset struct {
	debugLogger *log.Logger
	infoLogger  *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger
	logLevel    int
}

var lset loggerset

func init() {
	lset = loggerset{}
	lset.debugLogger = log.New(os.Stdout, "[DEBUG] ", log.Lshortfile|log.LstdFlags)
	lset.infoLogger = log.New(os.Stdout, "[INFO] ", log.Lshortfile|log.LstdFlags)
	lset.errorLogger = log.New(os.Stdout, "[ERROR] ", log.Lshortfile|log.LstdFlags)
	lset.fatalLogger = log.New(os.Stdout, "[FATAL] ", log.Lshortfile|log.LstdFlags)

	lset.logLevel = DebugLevel
}

// SetLevel is function to set log level
func SetLevel(level int) {
	lset.logLevel = level
}

// LogDebug print log content when logLevel <= Debug
func LogDebug(format string, v ...interface{}) {
	if lset.logLevel <= DebugLevel {
		lset.debugLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

// LogInfo print log content when logLevel <= Info
func LogInfo(format string, v ...interface{}) {
	if lset.logLevel <= InfoLevel {
		lset.infoLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

// LogError print log content when logLevel <= Error
func LogError(format string, v ...interface{}) {
	if lset.logLevel <= ErrorLevel {
		lset.errorLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

// FatalError print log content when logLevel <= Error
func LogFatal(format string, v ...interface{}) {
	if lset.logLevel <= FatalLevel {
		lset.fatalLogger.Output(2, fmt.Sprintf(format, v...))
		os.Exit(1)
	}
}
