package kleos

import (
	"encoding/json"
	"io"
	"strings"
	"sync"
)

// JSON field names for critical log data
const (
	FieldTimestamp = "ts"
	FieldMessage   = "msg"
	FieldLevel     = "level"
	FieldVerbosity = "v"
	FieldPkg       = "pkg"
	FieldSrc       = "src"
	FieldLine      = "line"
	FieldError     = "err"
)

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
type JSONOutput struct {
	sync.Mutex
	out     io.Writer
	encoder *json.Encoder
}

// NewJSONOutput creates a new log output that's meant to be used with the ELK stack.  Supports ECS
// fields for the standard fields.  See JSONOutput for details.
func NewJSONOutput(writer io.Writer) *JSONOutput {
	return &JSONOutput{
		out:     writer,
		encoder: json.NewEncoder(writer),
	}
}

func (w *JSONOutput) Write(m Message) error {
	if m.fields == nil {
		m.fields = make(Fields)
	}

	m.fields[FieldTimestamp] = m.when.UTC().Format(PaddedRFC3339Ms)

	if m.verbosity > 0 {
		m.fields[FieldLevel] = "debug"
		m.fields[FieldVerbosity] = m.verbosity
	} else if m.error == nil {
		m.fields[FieldLevel] = "info"
	} else {
		m.fields[FieldLevel] = "error"
	}

	// Write out the human-readable message
	msg := strings.TrimSpace(m.msg)
	if msg != "" {
		m.fields[FieldMessage] = msg
	}

	if m.file != "" {
		m.fields[FieldPkg] = m.pkg
		m.fields[FieldSrc] = m.file
		m.fields[FieldLine] = m.line
	}

	if m.error != nil {
		m.fields[FieldError] = m.error.Error()
	}

	// Applies any registered context variables to the fields
	contexts.Run(m.ctx, m.fields)

	w.Lock()
	defer w.Unlock()

	if err := w.encoder.Encode(m); err != nil {
		return err
	}

	return nil
}
