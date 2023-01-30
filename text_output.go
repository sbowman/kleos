package kleos

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// TextOutput is meant to output to stdout or stderr in black and white.
type TextOutput struct {
	sync.Mutex
	out io.Writer
}

// NewTextOutput creates a logger that writes human-readable plain text with no coloring.
func NewTextOutput(out io.Writer) *TextOutput {
	return &TextOutput{
		out: out,
	}
}

// Write the message out in plain text, but human-readable.
func (w *TextOutput) Write(m Message) error {
	w.Lock()
	defer w.Unlock()

	_, _ = fmt.Fprint(w.out, m.when.UTC().Format(PaddedRFC3339Ms))

	if m.verbosity > 0 {
		_, _ = fmt.Fprintf(w.out, " D%02d", m.verbosity)
	} else if m.error == nil {
		_, _ = fmt.Fprint(w.out, " INF")
	} else {
		_, _ = fmt.Fprint(w.out, " ERR")
	}

	// Write out the human-readable message
	msg := strings.TrimSpace(m.msg)
	if msg != "" {
		_, _ = fmt.Fprint(w.out, " ")
		_, _ = fmt.Fprint(w.out, msg)
	}

	// Where was the message logged?
	if m.file != "" {
		_, _ = fmt.Fprintf(w.out, " (%s/%s:%d)", m.pkg, m.file, m.line)
	}

	if m.error != nil {
		_, _ = fmt.Fprint(w.out, ", err=")
		_, _ = fmt.Fprint(w.out, strconv.Quote(m.error.Error()))
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
				_, _ = fmt.Fprint(w.out, ", ")
				_, _ = fmt.Fprint(w.out, k)
				_, _ = fmt.Fprint(w.out, "=")
				_, _ = fmt.Fprint(w.out, v)
			}
		}
	}

	_, _ = fmt.Fprintln(w.out)

	return nil
}
