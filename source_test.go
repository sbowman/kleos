package kleos_test

import (
	"bytes"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"

	"github.com/sbowman/kleos"
	"github.com/sbowman/tort"
)

func source() (string, string, int) {
	_, path, line, ok := runtime.Caller(1)
	if !ok {
		return "", "", -1
	}

	pkg := filepath.Base(filepath.Dir(path))
	file := filepath.Base(path)

	return pkg, file, line
}

func check(t *testing.T, output string, line int) {
	assert := tort.For(t)
	assert.String(output).Contains("(kleos/source_test.go:" + strconv.Itoa(line) + ")")
}

// Note:  if you change the internals of this test, you may need to update the "base+" values.
func TestLineNumbers(t *testing.T) {
	// Set a baseline
	_, _, base := source()

	var out bytes.Buffer
	kleos.SetOutput(kleos.NewTextOutput(&out))
	kleos.SetVerbosity(1)

	kleos.Log("Hello World") // +6
	check(t, out.String(), base+6)
	out.Reset()

	kleos.Debug("Hello World") // +10
	check(t, out.String(), base+10)
	out.Reset()

	kleos.V(1).With(kleos.Fields{
		"id":     "B8012423573231",
		"name":   "hello",
		"multi":  "taking space",
		"health": 97,
	}).Log("Hello World")
	check(t, out.String(), base+14)
	out.Reset()

	kleos.V(1).With(kleos.Fields{
		"id":     "B8012423573231",
		"name":   "hello",
		"multi":  "taking space",
		"health": 97,
	}).Debug("Hello World")
	check(t, out.String(), base+23)
	out.Reset()

	kleos.V(1).With(kleos.Fields{
		"id":     "B8012423573231",
		"name":   "hello",
		"multi":  "taking space",
		"health": 97,
	}).Info("Hello World")
	check(t, out.String(), base+32)
	out.Reset()
}

// Note:  if you change the internals of this test, you may need to update the "base+" values.
func TestLoggerLineNumbers(t *testing.T) {
	// Set a baseline
	_, _, base := source()

	var out bytes.Buffer
	kleos.SetOutput(kleos.NewTextOutput(&out))
	kleos.SetVerbosity(1)

	logger := kleos.NewLogger(1)
	logger.Printf("Hello world!")

	check(t, out.String(), base+7)
	out.Reset()
}
