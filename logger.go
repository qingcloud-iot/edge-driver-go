package edge_driver_go

import (
	"fmt"
	"os"
)

// Logger specifies logging API.
type Logger interface {
	// Debug logs any object in JSON format on debug level.
	Debug(a ...interface{})
	// Info logs any object in JSON format on info level.
	Info(a ...interface{})
	// Warn logs any object in JSON format on warning level.
	Warn(a ...interface{})
	// Error logs any object in JSON format on error level.
	Error(a ...interface{})
}
type logger struct {
}

func newLogger() Logger {
	return &logger{}
}
func (l *logger) Debug(a ...interface{}) {
	fmt.Fprintln(os.Stdout, a...)
}
func (l *logger) Info(a ...interface{}) {
	fmt.Fprintln(os.Stdout, a...)
}
func (l *logger) Warn(a ...interface{}) {
	fmt.Fprintln(os.Stdout, a...)
}
func (l *logger) Error(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
}
