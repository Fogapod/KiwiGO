package logger

import (
	"strconv"
	"testing"
)

func TestNewLogger(t *testing.T) {
	requiredLoggingLevel := TraceLevel

	l, err := NewLogger(requiredLoggingLevel)
	if err != nil {
		t.Fatalf("Failed to create logger instance with level %d:\n%s", requiredLoggingLevel, err)
	}

	if l.LoggingLevel != requiredLoggingLevel {
		t.Fatalf("Logging level does not match requested (%d instead of %d)", l.LoggingLevel, requiredLoggingLevel)
	}

	l, err = NewLogger(1.0)
	if err == nil {
		t.Fatal("Invalid type passed to NewLogger, but no error returned")
	}
}

func TestGetLogger(t *testing.T) {
	l := GetLogger()

	if l == nil {
		t.Fatal("Logger pointer is nil")
	}
}

func TestSetLoggingLevel(t *testing.T) {
	requiredLoggingLevel := DebugLevel
	requiredLoggingLevelString := "Debug"

	// test with int
	l, _ := NewLogger(TraceLevel)

	err := l.SetLoggingLevel(requiredLoggingLevel)
	if err != nil {
		t.Fatalf("Failed to set logging level to %d:\n%s", requiredLoggingLevel, err)
	}

	if l.LoggingLevel != requiredLoggingLevel {
		t.Fatalf("Logging level does not match requested (%d instead of %d)", l.LoggingLevel, requiredLoggingLevel)
	}

	// test with string
	l, _ = NewLogger(TraceLevel)

	err = l.SetLoggingLevel(requiredLoggingLevelString)
	if err != nil {
		t.Fatalf("Failed to set logging level to %s:\n%s", requiredLoggingLevelString, err)
	}

	if l.LoggingLevel != requiredLoggingLevel {
		t.Fatalf("Logging level does not match requested (%d instead of %d)", l.LoggingLevel, requiredLoggingLevel)
	}

	// test with string (number)
	l, _ = NewLogger(TraceLevel)

	err = l.SetLoggingLevel(strconv.Itoa(requiredLoggingLevel))
	if err != nil {
		t.Fatalf("Failed to set logging level with string %d:\n%s", requiredLoggingLevel, err)
	}

	if l.LoggingLevel != requiredLoggingLevel {
		t.Fatalf("Logging level does not match requested (%d instead of %d) (using string number)", l.LoggingLevel, requiredLoggingLevel)
	}

	// test with invalid string
	err = l.SetLoggingLevel("invalid value")
	if err == nil {
		t.Fatal("Invalid data passed to SetLoggingLevel, but no error was returned")
	}

	// test with strings
	for _, s := range []string{"Trace", "DEBUG", " info ", "warning", "error", "fatal"} {
		err = l.SetLoggingLevel(s)

		if err != nil {
			t.Fatalf("Failed to set logging level using key %s:\n%s", s, err)
		}
	}

	// test with invalid type
	err = l.SetLoggingLevel(1.0)
	if err == nil {
		t.Fatal("Invalid type passed to SetLoggingLevel, but no error was returned")
	}
}

func TestLogMessages(t *testing.T) {
	l, _ := NewLogger(TraceLevel)

	l.Fatal("[1/6] You should see this message")
	l.Error("[2/6] You should see this message")
	l.Warn("[3/6] You should see this message")
	l.Info("[4/6] You should see this message")
	l.Debug("[5/6] You should see this message")
	l.Trace("[6/6] You should see this message")

	l.SetLoggingLevel(FatalLevel)

	l.Fatal("[1/1] This should see this message")
	l.Error("You should not see this message")
	l.Warn("You should not see this message")
	l.Info("You should not see this message")
	l.Debug("You should not see this message")
	l.Trace("You should not see this message")
}
