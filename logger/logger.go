package logger

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	log *Logger
)

const (
	FatalLevel = iota
	DebugLevel
	InfoLevel
	TraceLevel
)

type Logger struct {
	LoggingLevel int
}

func init() {
	log, _ = NewLogger(FatalLevel)
}

func NewLogger(loggingLevel interface{}) (*Logger, error) {
	logger := &Logger{}
	if err := logger.SetLoggingLevel(loggingLevel); err != nil {
		return nil, err
	}

	return logger, nil
}

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
		switch strings.ToLower(strings.TrimSpace(value)) {
		case "fatal":
			return FatalLevel, nil
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

func GetLogger() *Logger {
	return log
}

// TODO: write to stdout instead of using Println
// TODO: write to file (optional?)
// TODO: better formatting for trailing unused arguments
func (l Logger) write(level int, prefix, postfix, msg string, args ...interface{}) {
	if level > l.LoggingLevel {
		return
	}

	now := time.Now().In(time.UTC)
	nowFormatted := now.Format(" 01-02 15:04:05 UTC ")

	msgFormatted := fmt.Sprintf(msg, args...)
	text := "\r"

	for _, line := range strings.Split(msgFormatted, "\n") {
		text += prefix + nowFormatted + line + postfix + "\n"
	}

	fmt.Print(text)
}

func (l Logger) Fatal(msg string, args ...interface{}) {
	l.write(0, "\u001b[31m[  !  ]", "\u001b[0m", msg, args...)
}

func (l Logger) Info(msg string, args ...interface{}) {
	l.write(1, "\u001b[32m[ INF ]", "\u001b[0m", msg, args...)
}

func (l Logger) Debug(msg string, args ...interface{}) {
	l.write(2, "\u001b[33m[DEBUG]", "\u001b[0m", msg, args...)
}

func (l Logger) Trace(msg string, args ...interface{}) {
	l.write(3, "[TRACE]", "", msg, args...)
}
