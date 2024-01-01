package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests if base implements basicModeler
func TestBasicModelerForBase(t *testing.T) {
	base := base{
		state:     Ready,
		ModelName: "BasicModelerTest",
		id:        123,
	}
	_, ok := any(base).(basicModeler)
	assert.True(t, ok, "base does NOT implement the basicModeler interface")
}
