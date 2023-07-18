package kleos_test

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/sbowman/kleos"
	"github.com/sbowman/tort"
)

func TestLogging(t *testing.T) {
	assert := tort.For(t)

	var out bytes.Buffer
	kleos.SetOutput(kleos.NewTextOutput(&out))
	kleos.SetVerbosity(1)

	kleos.V(1).With(kleos.Fields{
		"id":     "B8012423573231",
		"name":   "hello",
		"multi":  "taking space",
		"health": 97,
	}).Log("Hello World")

	output := out.String()
	assert.String(output).EndsWith("\n")
	assert.String(output).Contains("id=B8012423573231")
	assert.String(output).Contains("name=hello")
	assert.String(output).Contains("multi=\"taking space\"")
	assert.String(output).Contains("health=97")
	assert.String(output).Contains(" Hello World ")
	assert.String(output).Contains("D01")
}

func TestLevel(t *testing.T) {
	assert := tort.For(t)

	var out bytes.Buffer
	kleos.SetOutput(kleos.NewTextOutput(&out))
	kleos.SetVerbosity(1)

	kleos.With(kleos.Fields{
		"id":     "B8012423573231",
		"name":   "hello",
		"multi":  "taking space",
		"health": 97,
	}).Log("Hello World")

	output := out.String()
	assert.String(output).Contains("INF")
	out.Reset()

	kleos.V(1).With(kleos.Fields{
		"id":     "B8012423573231",
		"name":   "hello",
		"multi":  "taking space",
		"health": 97,
	}).Log("Hello World")

	output = out.String()
	assert.String(output).Contains("D01")
	out.Reset()

	kleos.With(kleos.Fields{
		"id":     "B8012423573231",
		"name":   "hello",
		"multi":  "taking space",
		"health": 97,
	}).Debug("Hello World")

	output = out.String()
	assert.String(output).Contains("D01")
	out.Reset()

	kleos.Error(fmt.Errorf("yikes")).With(kleos.Fields{
		"id":     "B8012423573231",
		"name":   "hello",
		"multi":  "taking space",
		"health": 97,
	}).Log("Hello World")

	output = out.String()
	assert.String(output).Contains("ERR")
	assert.String(output).Contains("err=yikes")
	out.Reset()
}

func TestLoggingAsync(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			time.Sleep(time.Millisecond)
			kleos.V(1).With(kleos.Fields{
				"id":     "B8012423573231",
				"name":   "hello",
				"health": 97,
			}).Log("Hello World")
		}()
	}
	wg.Wait()
}

func BenchmarkJSONLogging(b *testing.B) {
	b.ReportAllocs()
	kleos.SetOutput(kleos.NewJSONOutput(io.Discard))

	for n := 0; n < b.N; n++ {
		kleos.V(1).With(kleos.Fields{
			"id":     "B8012423573231",
			"name":   "NBC Sports",
			"health": 97,
		}).Log("Database is having serious issues related to connections")
	}
}
