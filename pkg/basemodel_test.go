package back

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests if base implements basicModeler
func TestBasicModelerForBase(t *testing.T) {
	base := Base{
		state:     Ready,
		ModelName: "BasicModelerTest",
		id:        123,
	}
	_, ok := any(base).(Basicmodeler)
	assert.True(t, ok, "base does NOT implement the basicModeler interface")
}
