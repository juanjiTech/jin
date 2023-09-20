package jin

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNameOfFunction(t *testing.T) {
	f := func() {}
	assert.Equal(t, "github.com/juanjiTech/jin.TestNameOfFunction.func1", nameOfFunction(f))
	assert.Equal(t, "nil", nameOfFunction(nil))
	assert.Equal(t, "nil", nameOfFunction(HandlerFunc(nil)))
}
