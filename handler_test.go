package jin

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFastInvokeWarpHandlerChain(t *testing.T) {
	assert.NotPanics(t, func() {
		fastInvokeWarpHandlerChain([]HandlerFunc{nil})
	})
}
