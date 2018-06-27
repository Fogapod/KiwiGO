package main

import (
	"fmt"
	"time"
)

// TODO: LoggingLevel constants?
type logger struct {
	LoggingLevel int
}

func getLogger() *logger {
	return &logger{0}
}

// TODO: write to stdout instead of using Println
// TODO: write to file (optional?)
// TODO: escape ANSII codes
// TODO: better formatting for trailing unused arguments
func (l logger) write(level int, prefix, msg string, args ...interface{}) {
	if level > l.LoggingLevel {
		return
	}

	now := time.Now().In(time.UTC)
	fmt.Println(fmt.Sprintf(prefix+now.Format(" 01-02 15:04:05 UTC ")+msg, args...))
}

// TODO terminal colours
func (l logger) Fatal(msg string, args ...interface{}) {
	l.write(0, "[  !  ]", msg, args...)
}

func (l logger) Info(msg string, args ...interface{}) {
	l.write(1, "[ INF ]", msg, args...)
}

func (l logger) Debug(msg string, args ...interface{}) {
	l.write(2, "[DEBUG]", msg, args...)
}

func (l logger) Trace(msg string, args ...interface{}) {
	l.write(3, "[TRACE]", msg, args...)
}
