package jin

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpHandlerFuncInvoker(t *testing.T) {
	var invoked bool
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		invoked = true
	})
	invoker := httpHandlerFuncInvoker(handler)
	_, err := invoker.Invoke([]interface{}{httptest.NewRecorder(), &http.Request{}})
	assert.NoError(t, err)
	assert.True(t, invoked)
}

func TestFastInvokeWarpHandler(t *testing.T) {
	t.Run("PanicNotFunc", func(t *testing.T) {
		assert.Panics(t, func() {
			fastInvokeWarpHandler("not a function")
		})
	})

	t.Run("AlreadyFastInvoker", func(t *testing.T) {
		handler := ContextInvoker(func(c *Context) {})
		warpedHandler := fastInvokeWarpHandler(handler)
		assert.Equal(t, reflect.ValueOf(handler), reflect.ValueOf(warpedHandler))
	})

	t.Run("HttpHandlerFunc", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
		warpedHandler := fastInvokeWarpHandler(handler)
		_, ok := warpedHandler.(httpHandlerFuncInvoker)
		assert.True(t, ok)
	})

	t.Run("GenericFunc", func(t *testing.T) {
		handler := func() {}
		warpedHandler := fastInvokeWarpHandler(handler)
		assert.Equal(t, reflect.ValueOf(handler).Pointer(), reflect.ValueOf(warpedHandler).Pointer())
	})
}

func TestFastInvokeWarpHandlerChain(t *testing.T) {
	chain := HandlersChain{
		func(c *Context) {},
		nil,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
	}

	fastInvokeWarpHandlerChain(chain)

	assert.NotNil(t, chain[0])
	assert.Nil(t, chain[1])
	assert.NotNil(t, chain[2])
	_, ok1 := chain[0].(ContextInvoker)
	assert.True(t, ok1)
	_, ok2 := chain[2].(httpHandlerFuncInvoker)
	assert.True(t, ok2)
}
