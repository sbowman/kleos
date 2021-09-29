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
	Out io.Writer
}

func NewTextOutput(out io.Writer) *TextOutput {
	return &TextOutput{
		Out:       out,
	}
}

func (w *TextOutput) Write(m Message, msg string, args ...interface{}) error {
	w.Lock()
	defer w.Unlock()

	_, _ = fmt.Fprint(w.Out, m.when.UTC().Format(PaddedRFC3339Ms))

	if m.debug {
		_, _ =fmt.Fprintf(w.Out, " DBG [%03d]", m.verbosity)
	} else if m.err == nil {
		_, _ = fmt.Fprint(w.Out, " INF")
	} else {
		_, _ = fmt.Fprint(w.Out, " ERR")
	}

	// Write out the human-readable message
	msg = strings.TrimSpace(msg)
	if msg != "" {
		_, _ = fmt.Fprint(w.Out, " ")
		_, _ =fmt.Fprintf(w.Out, msg, args...)
	}

	// Where was the message logged?
	if m.file != "" {
		_, _ = fmt.Fprintf(w.Out, " (%s/%s:%d)", m.pkg, m.file, m.line)
	}

	if m.err != nil {
		_, _ = fmt.Fprint(w.Out, ", error=")
		_, _ = fmt.Fprint(w.Out, strconv.Quote(m.err.Error()))
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
				_, _ = fmt.Fprint(w.Out, ", ")
				_, _ = fmt.Fprint(w.Out, k)
				_, _ = fmt.Fprint(w.Out, "=")
				_, _ = fmt.Fprint(w.Out, v)
			}
		}
	}

	fmt.Println()

	return nil
}
