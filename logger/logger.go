package logger

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	log *Logger
)

const (
	FatalLevel = iota
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

type Logger struct {
	LoggingLevel int
}

// Creates new global logger instance
func init() {
	log, _ = NewLogger(FatalLevel)
}

// Returns new logger instance
func NewLogger(loggingLevel interface{}) (*Logger, error) {
	logger := &Logger{}
	if err := logger.SetLoggingLevel(loggingLevel); err != nil {
		return nil, err
	}

	return logger, nil
}

// Set logging level
// Accepts string and int types as input
func (l *Logger) SetLoggingLevel(loggingLevel interface{}) error {
	value, err := parseLoggingLevel(loggingLevel)
	if err != nil {
		return err
	}

	l.LoggingLevel = value

	return nil
}

func parseLoggingLevel(input interface{}) (int, error) {
	switch input.(type) {
	case int:
		value, _ := input.(int)
		return value, nil
	case string:
		value, _ := input.(string)
		if number, err := strconv.Atoi(value); err == nil {
			return number, nil
		}

		switch strings.ToLower(strings.TrimSpace(value)) {
		case "fatal":
			return FatalLevel, nil
		case "error":
			return ErrorLevel, nil
		case "warning":
			return WarnLevel, nil
		case "info":
			return InfoLevel, nil
		case "debug":
			return DebugLevel, nil
		case "trace":
			return TraceLevel, nil
		default:
			return FatalLevel, errors.New("Invalid logging level name passed: " + value)
		}
	default:
		return FatalLevel, errors.New("Invalid logging level type passed")
	}
}

// Returns global logger instance. Init should be called before
func GetLogger() *Logger {
	return log
}

// TODO: write to stdout instead of using Print
// TODO: write to file (optional?)
// TODO: better formatting for trailing unused arguments
func (l Logger) write(level int, prefix, postfix, msg string, args ...interface{}) {
	if level > l.LoggingLevel {
		return
	}

	text := "\r"
	msgFormatted := fmt.Sprintf(msg, args...)
	nowFormatted := time.Now().In(time.UTC).Format(" 01-02 15:04:05 MST ")

	for _, line := range strings.Split(msgFormatted, "\n") {
		text += prefix + nowFormatted + line + postfix + "\n"
	}

	fmt.Print(text)
}

// TODO: make colours optional
// Send fatal message
func (l Logger) Fatal(msg string, args ...interface{}) {
	l.write(FatalLevel, "\u001b[31m[  !  ]", "\u001b[0m", msg, args...)
}

// Send error message
func (l Logger) Error(msg string, args ...interface{}) {
	l.write(ErrorLevel, "\u001b[31m[ ERR ]", "\u001b[0m", msg, args...)
}

// Send warning message
func (l Logger) Warn(msg string, args ...interface{}) {
	l.write(WarnLevel, "\u001b[31;1m[ WRN ]", "\u001b[0m", msg, args...)
}

// Send information message
func (l Logger) Info(msg string, args ...interface{}) {
	l.write(InfoLevel, "\u001b[32m[ INF ]", "\u001b[0m", msg, args...)
}

// Send debug message
func (l Logger) Debug(msg string, args ...interface{}) {
	l.write(DebugLevel, "\u001b[33m[DEBUG]", "\u001b[0m", msg, args...)
}

// Send trace message
func (l Logger) Trace(msg string, args ...interface{}) {
	l.write(TraceLevel, "[TRACE]", "", msg, args...)
}
