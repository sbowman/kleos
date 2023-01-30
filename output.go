package kleos

import (
	"fmt"
	"os"
)

const (
	// PaddedRFC3339Ms is the time format padded to three decimal places of ms.
	PaddedRFC3339Ms = "2006-01-02T15:04:05.000Z07:00"
)

// Output receives log messages; defaults to plain text output on stdout.
var output Writer = NewTextOutput(os.Stdout)

// Writer supports outputting log messages in various formats to various receivers, such as stdout
// or ELK.
type Writer interface {
	// Write a message to the output.  Messages should end in a carriage return.
	Write(m Message) error
}

// SetOutput changes the output writer.
func SetOutput(out Writer) {
	mutex.Lock()
	defer mutex.Unlock()

	output = out
}

// Output writes a nicely formatted message to the output device.
func (m Message) Output() {
	if err := output.Write(m); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to log message: %s", err)
		return
	}
}
