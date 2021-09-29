package kleos

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"
)

var EOL = []byte{'\n'}

// JSONOutput outputs in JSON format.  Meant for services like ELK or Splunk. Note that JSONOutput
// will overload these properties:
//
// * `ts` - when the message was generated
// * `host` - the hostname
// * `message`   - the plaintext log message
// * `level` - the log level, e.g. debug, info, error
// * `verbosity` - the verbosity of the debug message, when relevant
// * `pkg`   - the package in which this log message was generated
// * `src`   - the source file in which this log message was generated
// * `line`  - the source code line number that contains this log message
// * `error` - the error message formatted, if present
// * `code` - the error code, if the error supports it
//
type JSONOutput struct {
	sync.Mutex

	Host      string
	Timestamp string
	Out       io.Writer
}

// NewJSONOutput creates a new log output that's meant to be used with the ELK stack.  Supports ECS
// fields for the standard fields.  See JSONOutput for details.
func NewJSONOutput(host string, writer io.Writer) *JSONOutput {
	return &JSONOutput{
		Host:      host,
		Timestamp: "ts",
		Out:       writer,
	}
}

func (w *JSONOutput) Write(m Message, msg string, args ...interface{}) error {
	if m.fields == nil {
		m.fields = Fields{}
	}

	m.fields[w.Timestamp] = m.when.UTC().Format(PaddedRFC3339Ms)
	m.fields["host"] = w.Host

	if m.debug {
		m.fields["level"] = "debug"
		m.fields["verbosity"] = m.verbosity
	} else if m.err == nil {
		m.fields["level"] = "info"
	} else {
		m.fields["level"] = "error"
	}

	// Write out the human-readable message
	msg = strings.TrimSpace(msg)
	if msg != "" {
		m.fields["message"] = fmt.Sprintf(msg, args...)
	}

	if m.file != "" {
		m.fields["pkg"] = m.pkg
		m.fields["src"] = m.file
		m.fields["line"] = m.line
	}

	if m.err != nil {
		m.fields["error"] = m.err.Error()
		// TODO: Support error.code
	}

	// Applies any registered context variables to the fields
	contexts.Run(m.ctx, m.fields)

	w.Lock()
	defer w.Unlock()

	enc := json.NewEncoder(w.Out)
	if err := enc.Encode(m.fields); err != nil {
		return err
	}

	if _, err := w.Out.Write(EOL); err != nil {
		return err
	}

	return nil
}
