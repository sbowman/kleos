package kleos_test

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sbowman/kleos"
)

func TestLogging(t *testing.T) {
	assert := assert.New(t)

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
	assert.True(strings.HasSuffix(output, "\n"))
	assert.Contains(output, "id=B8012423573231")
	assert.Contains(output, "name=hello")
	assert.Contains(output, "multi=\"taking space\"")
	assert.Contains(output, "health=97")
	assert.Contains(output, " Hello World ")
	assert.Contains(output, "D01")
}

func TestLevel(t *testing.T) {
	assert := assert.New(t)

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
	assert.Contains(output, "INF")
	out.Reset()

	kleos.V(1).With(kleos.Fields{
		"id":     "B8012423573231",
		"name":   "hello",
		"multi":  "taking space",
		"health": 97,
	}).Log("Hello World")

	output = out.String()
	assert.Contains(output, "D01")
	out.Reset()

	kleos.With(kleos.Fields{
		"id":     "B8012423573231",
		"name":   "hello",
		"multi":  "taking space",
		"health": 97,
	}).Debug("Hello World")

	output = out.String()
	assert.Contains(output, "D01")
	out.Reset()

	kleos.Error(fmt.Errorf("yikes")).With(kleos.Fields{
		"id":     "B8012423573231",
		"name":   "hello",
		"multi":  "taking space",
		"health": 97,
	}).Log("Hello World")

	output = out.String()
	assert.Contains(output, "ERR")
	assert.Contains(output, "err=yikes")
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

func BenchmarkSimpleLogging(b *testing.B) {
	b.ReportAllocs()
	kleos.SetOutput(kleos.NewTextOutput(io.Discard))

	for n := 0; n < b.N; n++ {
		kleos.Log("Database is having serious issues related to connections")
	}
}

func BenchmarkNoSource(b *testing.B) {
	b.ReportAllocs()
	kleos.EnableSource(false)
	kleos.SetOutput(kleos.NewTextOutput(io.Discard))

	for n := 0; n < b.N; n++ {
		kleos.Log("Database is having serious issues related to connections")
	}
}
