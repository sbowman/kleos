package kleos

import (
	"context"
	"runtime"
	"time"
)

// Message carries details about the log message as function calls are made.  Note that
// messages aren't thread-safe, so don't pass them between goroutines (not sure why you'd
// do that).
type Message struct {
	k         *Kleos
	when      time.Time       // when was this message logged
	pkg       string          // in what package was the message generated
	file      string          // in what source code file was the message generated
	line      int             // what line of code generated the message
	ctx       context.Context // for parsing stored context values into the fields
	verbosity uint8           // verbosity level, 0-4
	msg       string          // the human-readable log message
	error     error           // include details about the error that generated this message
	fields    Fields          // any custom fields to include, typically as JSON output
	source    bool            // include the source file and line number?
	pc        []uintptr       // store the stacktrace
	skip      int             // how far back in the stacktrace to display source file and line number
	out       Writer
}

func generate(out Writer, source bool) Message {
	m := Message{
		when:   time.Now(),
		source: source,
		skip:   0,
		out:    out,
	}

	// A bit of extra effort so calling Source() repeatedly doesn't cost anything more
	if source {
		m.pc = make([]uintptr, 5)
		_ = runtime.Callers(3, m.pc)
	}

	return m
}

// Context records the context so that values stored in the context can be applied to the fields
// automatically on output.
func (m Message) Context(ctx context.Context) Message {
	m.ctx = ctx
	return m
}

// V applies a verbosity level to a debug message.
func (m Message) V(verbosity uint8) Message {
	m.verbosity = verbosity
	return m
}

// Error adds the error message as a field, "error", in the output.
func (m Message) Error(err error) Message {
	m.error = err
	return m
}

// With applies the given fields to the log message.
func (m Message) With(fields Fields) Message {
	m.fields = fields
	return m
}

// WithFields applies the given fields to the log message (deprecated).
func (m Message) WithFields(fields Fields) Message {
	m.fields = fields
	return m
}

// Source overrides the package, file, and line number of the log message.  Helpful for middleware.
func (m Message) Source(back int) Message {
	m.skip = back

	return m
}

// Debug generates a debug message.  Equivalent to `kleos.V(1).Log("This is a debug messsage!")`.
// If the Kleos verbosity is lower than the verbosity of the message, the message will not be
// output.
func (m Message) Debug(msg string) {
	m.msg = msg

	if m.verbosity < 1 {
		m.verbosity = 1
	}

	if m.verbosity > Verbosity() {
		return
	}

	m.Output()
}

// Log logs a message.  If the message has verbosity, it is logged as a debug message (or not logged
// if the Kleos verbosity setting isn't high enough).  If it has no verbosity but has errors, it
// is logged as an error message.  If it has no verbosity and no errors, it is logged as an info
// message.
func (m Message) Log(msg string) {
	m.msg = msg

	if m.verbosity > 0 && m.verbosity > Verbosity() {
		return
	}

	m.Output()
}

// Info logs a message.  Deprecated; use Log instead.
func (m Message) Info(msg string) {
	m.Log(msg)
}
