package kleos

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"
)

var EOL = []byte{'\n'}

// ELKOutput outputs in JSON format using the Elastic ECS format.  Meant for services like the ELK
// stack.  Note that ELKOutput will overload these properties:
//
// * `@timestamp` - when the message was generated
// * `message`   - the plaintext log message
// * `log.level` - the log level, e.g. debug, info, error
// * `verbosity` - the verbosity of the debug message, when relevant
// * `package.name`   - the package in which this log message was generated
// * `log.origin.file.name`   - the source file in which this log message was generated
// * `log.origin.file.line`  - the source code line number that contains this log message
// * `error.message` - the error message formatted, if present
// * `error.code` - the error code, if the error supports it
//
type ELKOutput struct {
	sync.Mutex

	Host string
	Out  io.Writer
}

// NewELKOutput creates a new log output that's meant to be used with the ELK stack.  Supports ECS
// fields for the standard fields.  See ELKOutput for details.
func NewELKOutput(host string, out io.Writer) *ELKOutput {
	return &ELKOutput{
		Host: host,
		Out:  out,
	}
}

func (w *ELKOutput) Write(m Message, msg string, args ...interface{}) error {
	if m.fields == nil {
		m.fields = Fields{}
	}

	m.fields["@timestamp"] = m.when.Format(PaddedRFC3339Ms)
	m.fields["host.name"] = w.Host

	if m.debug {
		m.fields["log.level"] = "debug"
		m.fields["verbosity"] = m.verbosity
	} else if m.err == nil {
		m.fields["log.level"] = "info"
	} else {
		m.fields["log.level"] = "error"
	}

	// Write out the human-readable message
	msg = strings.TrimSpace(msg)
	if msg != "" {
		m.fields["message"] = fmt.Sprintf(msg, args...)
	}

	if m.file != "" {
		m.fields["package.name"] = m.pkg
		m.fields["log.origin.file.name"] = m.file
		m.fields["log.origin.file.line"] = m.line
	}

	if m.err != nil {
		m.fields["error.message"] = m.err.Error()
		// TODO: Support error.code
	}

	// Applies any registered context variables to the fields
	contexts.Run(m.ctx, m.fields)

	enc := json.NewEncoder(w.Out)

	w.Lock()
	defer w.Unlock()

	if err := enc.Encode(m.fields); err != nil {
		return err
	}

	if _, err := w.Out.Write(EOL); err != nil {
		return err
	}

	return nil
}
