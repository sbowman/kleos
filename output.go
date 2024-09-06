package kleos

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const (
	// PaddedRFC3339Ms is the time format padded to three decimal places of ms.
	PaddedRFC3339Ms = "2006-01-02T15:04:05.000Z07:00"
)

// Writer supports outputting log messages in various formats to various receivers, such
// as stdout or ELK.
type Writer interface {
	// Write a message to the output.  Messages should end in a carriage return.
	Write(m Message) error
}

// SetOutput changes the output writer.
func SetOutput(out Writer) {
	local.SetOutput(out)
}

// SetOutput changes the output writer.
func (k *Kleos) SetOutput(out Writer) {
	k.Lock()
	defer k.Unlock()

	k.output = out
}

// Output writes a nicely formatted message to the output device.
func (m Message) Output() {
	if m.out == nil {
		return
	}

	if m.source && m.skip >= 0 && m.skip < len(m.pc) {
		frame, _ := runtime.CallersFrames(m.pc[m.skip : m.skip+1]).Next()
		_, file, line, ok := frame.PC, frame.File, frame.Line, frame.PC != 0
		if ok {
			m.pkg = filepath.Base(filepath.Dir(file))
			m.file = filepath.Base(file)
			m.line = line
		}
	}

	if err := m.out.Write(m); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to log message: %s", err)
		return
	}
}
