package logger

import (
	"fmt"
	"io"
	"os"
)

// Logger is the interface used throughout the compiler pipeline.
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}

// New returns a Logger that emits debug output only when debug=true.
func New(debug bool) Logger {
	if debug {
		return &stdLogger{w: os.Stderr}
	}
	return noop{}
}

// Noop returns a Logger that discards all output.
func Noop() Logger { return noop{} }

// stdLogger writes all levels to w.
type stdLogger struct{ w io.Writer }

func (l *stdLogger) Debug(args ...interface{}) {
	fmt.Fprintln(l.w, args...)
}

func (l *stdLogger) Debugf(format string, args ...interface{}) {
	fmt.Fprintf(l.w, format+"\n", args...)
}

func (l *stdLogger) Warn(args ...interface{}) {
	fmt.Fprint(l.w, "WARN: ")
	fmt.Fprintln(l.w, args...)
}

func (l *stdLogger) Warnf(format string, args ...interface{}) {
	fmt.Fprintf(l.w, "WARN: "+format+"\n", args...)
}

func (l *stdLogger) Error(args ...interface{}) {
	fmt.Fprint(l.w, "ERROR: ")
	fmt.Fprintln(l.w, args...)
}

func (l *stdLogger) Errorf(format string, args ...interface{}) {
	fmt.Fprintf(l.w, "ERROR: "+format+"\n", args...)
}

// noop discards everything.
type noop struct{}

func (noop) Debug(args ...interface{})            {}
func (noop) Debugf(string, ...interface{})        {}
func (noop) Warn(args ...interface{})             {}
func (noop) Warnf(string, ...interface{})         {}
func (noop) Error(args ...interface{})            {}
func (noop) Errorf(string, ...interface{})        {}
