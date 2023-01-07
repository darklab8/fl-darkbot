package consoler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGettingOutput(t *testing.T) {
	assert.Contains(t, Consoler{}.New("ping").Execute().GetResult(), "Pong!")
}
