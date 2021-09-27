package kleos

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

const (
	// PaddedRFC3339Ms is the time format padded to three decimal places of ms.
	PaddedRFC3339Ms = "2006-01-02T15:04:05.000Z07:00"
)

// Colorized output
var (
	timestamp = color.New(color.FgCyan, color.Faint)
	info      = color.New(color.FgGreen)
	err       = color.New(color.FgRed)
	debug     = color.New(color.FgMagenta)
	message   = color.New(color.FgHiWhite)
	field     = color.New(color.FgWhite, color.Faint)
	location  = color.New(color.FgCyan)
)

// Output writes a nicely formatted message to the output device.
func (m Message) Output(out io.Writer, msg string, args ...interface{}) {
	var send bytes.Buffer

	_, _ = timestamp.Fprint(&send, m.when.Format(PaddedRFC3339Ms))

	if m.debug {
		_, _ = debug.Fprintf(&send, " DBG [%03d]", m.verbosity)
	} else if m.err == nil {
		_, _ = info.Fprint(&send, " INF")
	} else {
		_, _ = err.Fprint(&send, " ERR")
	}

	// Write out the human-readable message
	msg = strings.TrimSpace(msg)
	if msg != "" {
		_, _ = fmt.Fprint(&send, " ")
		_, _ = message.Fprintf(&send, msg, args...)
	}

	// Where was the message logged?
	if m.file != "" {
		_, _ = location.Fprintf(&send, " (%s/%s:%d)", m.pkg, m.file, m.line)
	}

	if m.err != nil {
		_, _ = field.Fprint(&send, ", error=")
		_, _ = field.Fprint(&send, strconv.Quote(m.err.Error()))
	}

	if m.fields == nil {
		m.fields = Fields{}
	}

	// Applies any registered context variables to the fields
	contexts.Run(m.ctx, m.fields)

	if len(m.fields) > 0 {
		// Write the fields in alphabetical order
		var fields []string
		for field := range m.fields {
			fields = append(fields, field)
		}
		sort.Strings(fields)

		for _, k := range fields {
			v := encode(m.fields[k])

			if v != "" {
				_, _ = field.Fprint(&send, ", ")
				_, _ = field.Fprint(&send, k)
				_, _ = field.Fprint(&send, "=")
				_, _ = field.Fprint(&send, v)
			}
		}
	}

	// Ensures log messages are published as a single unit, so threads don't step on each other
	_, _ = fmt.Fprintln(out, send.String())
}
