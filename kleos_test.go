package kleos_test

import (
	"bytes"
	"io/ioutil"
	"strings"
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

	kleos.V(1).With(kleos.Fields{
		"id":     "B8012423573231",
		"name":   "hello",
		"health": 97,
	}).Log("Hello World")

	output := out.String()
	assert.IsTrue(strings.HasSuffix(output, "\n"))
	// TODO: other
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
	kleos.SetOutput(kleos.NewJSONOutput(ioutil.Discard))

	for n := 0; n < b.N; n++ {
		kleos.V(1).With(kleos.Fields{
			"id":     "B8012423573231",
			"name":   "NBC Sports",
			"health": 97,
		}).Log("Database is having serious issues related to connections")
	}
}
