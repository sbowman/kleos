package kleos_test

import (
	"bytes"
	"testing"

	"github.com/sbowman/kleos"
	"github.com/sbowman/tort"
)

// TODO: test that verbosity works right

func TestVerbosity(t *testing.T) {
	assert := tort.For(t)

	var out bytes.Buffer
	kleos.SetOutput(kleos.NewTextOutput(&out))
	kleos.SetVerbosity(2)

	kleos.Log("Hello World")

	output := out.String()
	assert.String(output).NotContains("D01")
	assert.String(output).Contains("INF")
	assert.String(output).Contains(" Hello World ")

	out.Reset()

	kleos.V(2).Log("Hello World")

	output = out.String()
	assert.String(output).Contains("D02")
	assert.String(output).NotContains("INF")
	assert.String(output).Contains(" Hello World ")

	out.Reset()

	kleos.V(3).Log("Hello World")

	output = out.String()
	assert.String(output).NotContains("D02")
	assert.String(output).NotContains("INF")
	assert.String(output).NotContains(" Hello World ")

	out.Reset()
}
