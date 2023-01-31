package kleos

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"

	"github.com/fatih/color"
)

// ColorOutput is meant to output to stdout or stderr with color.
type ColorOutput struct {
	sync.Mutex
	out io.Writer

	timestamp, info, err, debug, message, field, location *color.Color
}

// NewColorOutput creates a color output writer meant for stdout or stderr.
func NewColorOutput(out io.Writer) *ColorOutput {
	return &ColorOutput{
		out:       out,
		timestamp: color.New(color.FgCyan, color.Faint),
		info:      color.New(color.FgGreen),
		err:       color.New(color.FgRed),
		debug:     color.New(color.FgMagenta),
		message:   color.New(color.FgHiWhite),
		field:     color.New(color.FgWhite, color.Faint),
		location:  color.New(color.FgCyan),
	}
}

// Write the message to the color output writer.
func (w *ColorOutput) Write(m Message) error {
	w.Lock()
	defer w.Unlock()

	_, _ = w.timestamp.Fprint(w.out, m.when.UTC().Format(PaddedRFC3339Ms))

	if m.verbosity > 0 {
		_, _ = w.debug.Fprintf(w.out, " D%02d", m.verbosity)
	} else if m.error == nil {
		_, _ = w.info.Fprint(w.out, " INF")
	} else {
		_, _ = w.err.Fprint(w.out, " ERR")
	}

	// Write out the human-readable message
	msg := strings.TrimSpace(m.msg)
	if msg != "" {
		_, _ = fmt.Fprint(w.out, " ")
		_, _ = w.message.Fprint(w.out, msg)
	}

	// Where was the message logged?
	if m.file != "" {
		_, _ = w.location.Fprintf(w.out, " (%s/%s:%d)", m.pkg, m.file, m.line)
	}

	if m.error != nil {
		_, _ = w.field.Fprint(w.out, ", err=")
		_, _ = w.field.Fprint(w.out, encode(m.error.Error()))
	}

	if m.fields == nil {
		m.fields = make(Fields)
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
				_, _ = w.field.Fprint(w.out, ", ")
				_, _ = w.field.Fprint(w.out, k)
				_, _ = w.field.Fprint(w.out, "=")
				_, _ = w.field.Fprint(w.out, v)
			}
		}
	}

	_, _ = fmt.Fprintln(w.out)

	return nil
}
