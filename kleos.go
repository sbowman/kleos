// Package kleos is a simple, efficient logging service.  In development mode it outputs colorized
// log files for easy reading and review.  In production it produces JSON log files compatible with
// systems like ELK.
package kleos

import (
	"context"
	"os"
	"sync"
)

// Used for our "global" logger.
var local = New()

// Kleos represents a logger.  There's an internal logger instance that's created and used
// for the global functions such as `Log` or `With`, though you may create your own Kleos
// instance if you like.
type Kleos struct {
	sync.RWMutex

	output        Writer
	includeSource bool
}

// New creates a new logging instance.  Typically there's no need to do this unless you're
// being strict about no global variables.
func New() *Kleos {
	return &Kleos{
		output:        NewTextOutput(os.Stdout),
		includeSource: true,
	}
}

// EnableSource enables or disables reporting the source file and line number of the log
// message.
func (k *Kleos) EnableSource(enabled bool) {
	k.Lock()
	defer k.Unlock()

	k.includeSource = enabled
}

// EnableSource enables or disables reporting the source file and line number of the log
// message.
func EnableSource(enabled bool) {
	local.EnableSource(enabled)
}

// Context records the context so that values stored in the context can be applied to the
// fields automatically on output.
func (k *Kleos) Context(ctx context.Context) Message {
	return generate(k.output, k.includeSource).Context(ctx)
}

// V applies a verbosity level to a debug message.
func (k *Kleos) V(verbosity uint8) Message {
	return generate(k.output, k.includeSource).V(verbosity)
}

// Error adds the error message as a field, "source", in the output.
func (k *Kleos) Error(err error) Message {
	return generate(k.output, k.includeSource).Error(err)
}

// With applies the given fields to the log message.
func (k *Kleos) With(fields Fields) Message {
	return generate(k.output, k.includeSource).With(fields)
}

// WithFields applies the given fields to the log message (deprecated).
func (k *Kleos) WithFields(fields Fields) Message {
	return generate(k.output, k.includeSource).With(fields)
}

// Source overrides the package, file, and line number of the log message.  Helpful for
// middleware.
func (k *Kleos) Source(back int) Message {
	return generate(k.output, k.includeSource).Source(back)
}

// Debug generates a debug message.  Equivalent to `kleos.V(1).Log("This is a debug
// messsage!")`.  If the Kleos verbosity is lower than the verbosity of the message, the
// message will not be output.  Should use `V().Log()` instead.
func (k *Kleos) Debug(msg string) {
	generate(k.output, k.includeSource).Debug(msg)
}

// Log logs a message.  If the message has verbosity, it is logged as a debug message (or
// not logged if the Kleos verbosity setting isn't high enough).  If it has no verbosity
// but has errors, it is logged as an error message.  If it has no verbosity and no
// errors, it is logged as an info message.
func (k *Kleos) Log(msg string) {
	generate(k.output, k.includeSource).Log(msg)
}

// Info logs a message.  Deprecated; use Log instead.
func (k *Kleos) Info(msg string) {
	generate(k.output, k.includeSource).Log(msg)
}

// TODO: create a Logger struct and use that for the global logger.
// TODO: put the mutex in the logger

const pkgoffset = 1

// Context records the context so that values stored in the context can be applied to the
// fields automatically on output.
func Context(ctx context.Context) Message {
	return local.Source(pkgoffset).Context(ctx)
}

// V applies a verbosity level to a debug message.
func V(verbosity uint8) Message {
	return local.Source(pkgoffset).V(verbosity)
}

// Error adds the error message as a field, "source", in the output.
func Error(err error) Message {
	return local.Source(pkgoffset).Error(err)
}

// With applies the given fields to the log message.
func With(fields Fields) Message {
	return local.Source(pkgoffset).With(fields)
}

// WithFields applies the given fields to the log message (deprecated).
func WithFields(fields Fields) Message {
	return local.Source(pkgoffset).With(fields)
}

// Source overrides the package, file, and line number of the log message.  Helpful for
// middleware.
func Source(back int) Message {
	return local.Source(pkgoffset).Source(back + 1)
}

// Debug generates a debug message.  Equivalent to `kleos.V(1).Log("This is a debug
// messsage!")`.  If the Kleos verbosity is lower than the verbosity of the message, the
// message will not be output.  Should use `V().Log()` instead.
func Debug(msg string) {
	local.Source(pkgoffset).Debug(msg)
}

// Log logs a message.  If the message has verbosity, it is logged as a debug message (or
// not logged if the Kleos verbosity setting isn't high enough).  If it has no verbosity
// but has errors, it is logged as an error message.  If it has no verbosity and no
// errors, it is logged as an info message.
func Log(msg string) {
	local.Source(pkgoffset).Log(msg)
}

// Info logs a message.  Deprecated; use Log instead.
func Info(msg string) {
	local.Source(pkgoffset).Log(msg)
}
