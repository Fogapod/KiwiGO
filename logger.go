package main

import (
	"fmt"
	"strings"
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
// TODO: better formatting for trailing unused arguments
func (l logger) write(level int, prefix, postfix, msg string, args ...interface{}) {
	if level > l.LoggingLevel {
		return
	}

	now := time.Now().In(time.UTC)
	nowFormatted := now.Format(" 01-02 15:04:05 UTC ")

	msgFormatted := fmt.Sprintf(msg, args...)
	text := ""

	for _, line := range strings.Split(msgFormatted, "\n") {
		text += prefix + nowFormatted + line + postfix + "\n"
	}

	fmt.Print(text)
}

func (l logger) Fatal(msg string, args ...interface{}) {
	l.write(0, "\u001b[31m[  !  ]", "\u001b[0m", msg, args...)
}

func (l logger) Info(msg string, args ...interface{}) {
	l.write(1, "\u001b[32m[ INF ]", "\u001b[0m", msg, args...)
}

func (l logger) Debug(msg string, args ...interface{}) {
	l.write(2, "\u001b[33m[DEBUG]", "\u001b[0m", msg, args...)
}

func (l logger) Trace(msg string, args ...interface{}) {
	l.write(3, "[TRACE]", "", msg, args...)
}
