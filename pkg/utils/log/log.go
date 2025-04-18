package logging

import (
	"fmt"
	"strings"
)

type Level int

const (
	TraceLevel Level = iota
	DebugLevel
	InfoLevel
	ErrorLevel
	NoneLevel
)

func (l Level) String() string {
	return []string{"TRACE", "DEBUG", "INFO", "ERROR"}[l]
}

var level Level = ErrorLevel
var stack []string

func SetLevel(l Level) {
	level = l
}

// Trace logs a message at the trace severity level.
func Trace(args ...any) {
	log(TraceLevel, args...)
}

// Debug logs a message at the debug severity level.
func Debug(args ...any) {
	log(DebugLevel, args...)
}

func Debugf(format string, args ...any) {
	Debug(fmt.Sprintf(format, args...))
}

// Info logs a message at the info severity level.
func Info(args ...any) {
	log(InfoLevel, args...)
}

// Infof logs a message at the info severity level.
func Infof(format string, args ...any) {
	Info(fmt.Sprintf(format, args...))
}

// Error logs a message at the error severity level.
//
// This should be only called when an unhandled error from a
// third party package occurs.
func Error(args ...any) {
	log(ErrorLevel, args...)
	logStack()
}

// Errorf logs a message at the error severity level.
func Errorf(format string, args ...any) {
	Error(fmt.Sprintf(format, args...))
}

// Fn logs the entry and exit of a function with the given name.
// the logging is done at the trace severity level.
//
// Example:
//
//	func foo() {
//		defer logging.Fn("foo")()
//		// ...
//	}
//
// Output:
//
//	[TRACE] entering foo
//	[TRACE] exiting foo
func Fn(name string) func() {
	Trace("enter ", name)
	stack = append(stack, name)
	return func() {
		Trace("leave ", name)
		stack = stack[:len(stack)-1]
	}
}

func logStack() {
	fmt.Println("stack:")
	for _, f := range stack {
		fmt.Println("  ", f)
	}
}

func log(l Level, args ...any) {
	if l < level {
		return
	}

	prefix := fmt.Sprintf("[%s]", l.String())
	fmt.Println(prefix, join(args...))
}

func join(args ...any) string {
	s := make([]string, len(args))

	for i, a := range args {
		s[i] = fmt.Sprint(a)
	}

	return strings.Join(s, " ")
}
