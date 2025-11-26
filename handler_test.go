package jin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFastInvokeWarpHandlerChain(t *testing.T) {
	assert.NotPanics(t, func() {
		fastInvokeWarpHandlerChain([]HandlerFunc{nil})
	})
}
