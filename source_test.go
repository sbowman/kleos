package kleos_test

import (
	"bytes"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/sbowman/kleos"
	"github.com/stretchr/testify/assert"
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
	assert := assert.New(t)
	assert.Contains(output, fmt.Sprintf("(kleos/source_test.go:%d)", line))
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
	}).Log("Hello World 1")
	check(t, out.String(), base+14)
	out.Reset()

	kleos.With(kleos.Fields{
		"id":     "B8012423573231",
		"name":   "hello",
		"multi":  "taking space",
		"health": 97,
	}).V(2).Debug("Hello World 2")
	assert.Empty(t, out.String())
	out.Reset()

	kleos.V(1).With(kleos.Fields{
		"id":     "B8012423573231",
		"name":   "hello",
		"multi":  "taking space",
		"health": 97,
	}).Info("Hello World 3")
	check(t, out.String(), base+32)
	out.Reset()

	kleos.Info("Hello World")
	check(t, out.String(), base+41)
	out.Reset()
}

func TestKleosLineNumbers(t *testing.T) {
	// Set a baseline
	_, _, base := source()

	var out bytes.Buffer
	log := kleos.New()
	log.SetOutput(kleos.NewTextOutput(&out))
	log.SetVerbosity(1)

	log.Log("Hello World") // +6
	check(t, out.String(), base+7)
	out.Reset()

	log.Debug("Hello World") // +10
	check(t, out.String(), base+11)
	out.Reset()

	log.V(1).With(kleos.Fields{
		"id":     "B8012423573231",
		"name":   "hello",
		"multi":  "taking space",
		"health": 97,
	}).Log("Hello World 1")
	check(t, out.String(), base+15)
	out.Reset()

	log.With(kleos.Fields{
		"id":     "B8012423573231",
		"name":   "hello",
		"multi":  "taking space",
		"health": 97,
	}).V(2).Debug("Hello World 2")
	assert.Empty(t, out.String())
	out.Reset()

	log.V(1).With(kleos.Fields{
		"id":     "B8012423573231",
		"name":   "hello",
		"multi":  "taking space",
		"health": 97,
	}).Info("Hello World 3")
	check(t, out.String(), base+33)
	out.Reset()

	log.Info("Hello World")
	check(t, out.String(), base+42)
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

// Note:  if you change the internals of this test, you may need to update the "base+" values.
func TestSourceLineNumbers(t *testing.T) {
	// Set a baseline
	_, _, base := source()

	var out bytes.Buffer
	kleos.SetOutput(kleos.NewTextOutput(&out))

	deepCall()
	check(t, out.String(), base+5)
	out.Reset()
}

func deepCall() {
	kleos.Source(1).Debug("Hello World") // +10

}
