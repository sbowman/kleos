// Package kleos is a simple, efficient logging service.  In development mode it outputs colorized
// log files for easy reading and review.  In production it produces JSON log files compatible with
// systems like ELK.
package kleos

import (
	"context"
	"sync"
)

var mutex sync.RWMutex

// Context records the context so that values stored in the context can be applied to the fields
// automatically on output.
func Context(ctx context.Context) Message {
	return generate().Context(ctx)
}

// V applies a verbosity level to a debug message.
func V(verbosity uint8) Message {
	return generate().V(verbosity)
}

// Error adds the error message as a field, "source", in the output.
func Error(err error) Message {
	return generate().Error(err)
}

// With applies the given fields to the log message.
func (m Message) With(fields Fields) Message {
	return generate().With(fields)
}

// Debug generates a debug message.  Equivalent to `kleos.V(1).Log("This is a debug messsage!")`.
// If the Kleos verbosity is lower than the verbosity of the message, the message will not be
// output.
func Debug(msg string) {
	generate().Debug(msg)
}

// Log logs a message.  If the message has verbosity, it is logged as a debug message (or not logged
// if the Kleos verbosity setting isn't high enough).  If it has no verbosity but has errors, it
// is logged as an error message.  If it has no verbosity and no errors, it is logged as an info
// message.
func Log(msg string) {
	generate().Log(msg)
}
