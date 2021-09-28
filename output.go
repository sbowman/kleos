package kleos

import (
	"fmt"
	"os"
)

const (
	// PaddedRFC3339Ms is the time format padded to three decimal places of ms.
	PaddedRFC3339Ms = "2006-01-02T15:04:05.000Z07:00"
)

var (
	// Output receives log messages; defaults to a ColorOutput on stdout.
	output Writer
)

type Writer interface{
	// Write a message to the output.  Messages should end in a carriage return.
	Write(m Message, msg string, args ...interface{}) error
}

func init() {
	output = NewColorOutput(os.Stdout)
}

// Output writes a nicely formatted message to the output device.
func (m Message) Output(msg string, args ...interface{}) {
	if err := output.Write(m, msg, args...); err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "Unable to log message: %s", err)
	}
}
