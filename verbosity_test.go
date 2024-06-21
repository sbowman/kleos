package kleos_test

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/sbowman/kleos"
)

func TestVerbosity(t *testing.T) {
	assert := assert.New(t)

	var out bytes.Buffer
	kleos.SetOutput(kleos.NewTextOutput(&out))
	kleos.SetVerbosity(2)

	kleos.Log("Hello World")

	output := out.String()
	assert.NotContains(output, "D01")
	assert.Contains(output, "INF")
	assert.Contains(output, " Hello World ")

	out.Reset()

	kleos.V(2).Log("Hello World")

	output = out.String()
	assert.Contains(output, "D02")
	assert.NotContains(output, "INF")
	assert.Contains(output, " Hello World ")

	out.Reset()

	kleos.V(3).Log("Hello World")

	output = out.String()
	assert.NotContains(output, "D02")
	assert.NotContains(output, "INF")
	assert.NotContains(output, " Hello World ")

	out.Reset()
}
