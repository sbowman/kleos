package kleos

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/fatih/color"
)

// ColorOutput is meant to output to stdout or stderr with color.
type ColorOutput struct {
	sync.Mutex
	Out io.Writer

	timestamp, info, err, debug, message, field, location *color.Color
}

func NewColorOutput(out io.Writer) *ColorOutput {
	return &ColorOutput{
		Out:       out,
		timestamp: color.New(color.FgCyan, color.Faint),
		info:      color.New(color.FgGreen),
		err:       color.New(color.FgRed),
		debug:     color.New(color.FgMagenta),
		message:   color.New(color.FgHiWhite),
		field:     color.New(color.FgWhite, color.Faint),
		location:  color.New(color.FgCyan),
	}
}

func (w *ColorOutput) Write(m Message, msg string, args ...interface{}) error {
	w.Lock()
	defer w.Unlock()

	_, _ = w.timestamp.Fprint(w.Out, m.when.UTC().Format(PaddedRFC3339Ms))

	if m.debug {
		_, _ = w.debug.Fprintf(w.Out, " DBG [%03d]", m.verbosity)
	} else if m.err == nil {
		_, _ = w.info.Fprint(w.Out, " INF")
	} else {
		_, _ = w.err.Fprint(w.Out, " ERR")
	}

	// Write out the human-readable message
	msg = strings.TrimSpace(msg)
	if msg != "" {
		_, _ = fmt.Fprint(w.Out, " ")
		_, _ = w.message.Fprintf(w.Out, msg, args...)
	}

	// Where was the message logged?
	if m.file != "" {
		_, _ = w.location.Fprintf(w.Out, " (%s/%s:%d)", m.pkg, m.file, m.line)
	}

	if m.err != nil {
		_, _ = w.field.Fprint(w.Out, ", error=")
		_, _ = w.field.Fprint(w.Out, strconv.Quote(m.err.Error()))
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
				_, _ = w.field.Fprint(w.Out, ", ")
				_, _ = w.field.Fprint(w.Out, k)
				_, _ = w.field.Fprint(w.Out, "=")
				_, _ = w.field.Fprint(w.Out, v)
			}
		}
	}

	fmt.Println()

	return nil
}
