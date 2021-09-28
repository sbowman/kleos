// Package kleos is a simple, efficient logging service.  In development mode it outputs colorized
// log files for easy reading and review.  In production it produces JSON log files compatible with
// systems like ELK.
package kleos

import (
	"context"
	"path/filepath"
	"runtime"
	"time"
)

var (
	// Verbosity is the default debug level to output.  If less than zero, no debug
	// messages are output.
	verbosity = -1
)

// Fields holds details about the log message.
type Fields map[string]interface{}

// Message carries details about the log message as function calls are made.
type Message struct {
	when      time.Time
	ctx       context.Context
	pkg       string
	file      string
	line      int
	verbosity int  // debug messages
	debug     bool // is this a debug message?
	err       error
	fields    Fields
}

// WithFields applies the given fields to the log message.
func WithFields(fields Fields) Message {
	return Message{
		when:   time.Now(),
		fields: fields,
	}.FileAndLineNumber()
}

// WithContext associates the context with the log message, to use the user and
// request ID in the logs.
func WithContext(ctx context.Context) Message {
	return Message{
		when:    time.Now(),
		ctx: ctx,
	}.FileAndLineNumber()
}

// Verbosity returns the current verbosity
func Verbosity() int {
	return verbosity
}

// V applies a verbosity level to a debug message.
func V(verbosity int) Message {
	return Message{
		when:      time.Now(),
		verbosity: verbosity,
	}.FileAndLineNumber()
}

// Error adds the error message as a field, "error", in the output.
func Error(err error) Message {
	return Message{
		when: time.Now(),
		err:  err,
	}.FileAndLineNumber()
}

// Context references the current context in the log message, recording the request ID if present.
func Context(ctx context.Context) Message {
	return Message{
		when: time.Now(),
		ctx:  ctx,
	}.FileAndLineNumber()
}

// FileAndLineNumber appends the filename and line number information to the
// log message (where the log was implemented).
func (m Message) FileAndLineNumber() Message {
	return m.FileAndLineNumberBack(3)
}

// FileAndLineNumberBack appends the filename and line number information to the
// log message (where the log was implemented).
func (m Message) FileAndLineNumberBack(back int) Message {
	_, file, line, ok := runtime.Caller(back)
	if !ok {
		return m
	}

	m.pkg = filepath.Base(filepath.Dir(file))
	m.file = filepath.Base(file)
	m.line = line

	return m
}

// WithFields applies the given fields to the log message.
func (m Message) WithFields(fields Fields) Message {
	m.fields = fields
	return m
}

// WithContext associates the context with the log message, to use the user and
// request ID in the logs.
func (m Message) WithContext(ctx context.Context) Message {
	m.ctx = ctx
	return m
}

// Verbosity applies a verbosity level to a debug message.
// func (m Message) Verbosity(verbosity int) Message {
// 	m.verbosity = verbosity
// 	return m
// }

// V applies a verbosity level to a debug message.
func (m Message) V(verbosity int) Message {
	m.verbosity = verbosity
	return m
}

// Error adds the error message as a field, "error", in the output.
func (m Message) Error(err error) Message {
	m.err = err
	return m
}

// Context references the current context in the log message, recording the request ID if present.
func (m Message) Context(ctx context.Context) Message {
	m.ctx = ctx
	return m
}

// Info writes an info-level message.  Info messages are always written.
func (m Message) Info(msg string) {
	m.Output(msg)
}

// Infof writes a formatted info-level message, using printf formatting.
func (m Message) Infof(msg string, args ...interface{}) {
	m.Output(msg, args...)
}

// Debug writes an debug-level message.  Debug messages are written when the
// verbosity is lower than the logging verbosity.
func (m Message) Debug(msg string) {
	if verbosity < m.verbosity {
		return
	}

	m.debug = true
	m.Output(msg)
}

// Info writes an info-level message.  Info messages are always written.
func Info(msg string) {
	Message{
		when: time.Now(),
	}.FileAndLineNumber().Output(msg)
}

// Debug writes an debug-level message.  Debug messages are written when the
// verbosity is lower than the logging verbosity.
func Debug(msg string) {
	if verbosity < 0 {
		return
	}

	m := Message{
		when:  time.Now(),
		debug: true,
	}

	m.FileAndLineNumber().Output(msg)
}

// SetOutput changes the output writer.
func SetOutput(out Writer) {
	output = out
}

// SetVerbosity sets the debug verbosity level.  May be set on the fly.
func SetVerbosity(level int) {
	verbosity = level
}
