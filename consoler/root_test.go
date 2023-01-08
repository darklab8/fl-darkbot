package consoler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGettingOutput(t *testing.T) {
	assert.Contains(t, Consoler{}.New("ping").Execute().GetResult(), "Pong!")
}

func TestGrabStdout(t *testing.T) {
	c := Consoler{}.New("ping --help")
	result := c.Execute().GetResult()

	assert.Contains(t, result, "\nFlags:\n  -h, --help   ")
}
