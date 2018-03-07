package common

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestFloatArr_Append(t *testing.T) {
	fa := NewFloatArr()
	fa.Append(1).Append(2)

	assert.InDelta(t, 1, fa.At(0), 1e-9)
	assert.InDelta(t, 2, fa.At(1), 1e-9)
}
