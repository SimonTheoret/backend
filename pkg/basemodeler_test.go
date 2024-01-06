package back

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests if base implements the modeler interface
func TestBaseImplementModelerInterface(t *testing.T) {
	sender := HttpModel{dest: "anything"}
	base := NewBase{"test", sender}
	_, ok := any(&base).(modeler) //Must be &h because an interface is always a pointer type
	assert.True(t, ok, "base does NOT implement the modeler interface")
}
