package jin

import (
	"fmt"
	"github.com/juanjiTech/inject/v2"
	"net/http"
	"reflect"
)

// HandlerFunc defines the handler used by Jin middleware as return value.
type HandlerFunc interface{}

// HandlersChain defines a HandlerFunc slice.
type HandlersChain []HandlerFunc

// Last returns the last handler in the chain. i.e. the last handler is the main one.
func (c HandlersChain) Last() HandlerFunc {
	if length := len(c); length > 0 {
		return c[length-1]
	}
	return nil
}

var _ inject.FastInvoker = (*ContextInvoker)(nil)

// ContextInvoker is an inject.FastInvoker implementation of `func(Context)`.
type ContextInvoker func(ctx *Context)

func (invoke ContextInvoker) Invoke(args []interface{}) ([]reflect.Value, error) {
	invoke(args[0].(*Context))
	return nil, nil
}

// httpHandlerFuncInvoker is an inject.FastInvoker implementation of
// `func(http.ResponseWriter, *http.Request)`.
type httpHandlerFuncInvoker func(http.ResponseWriter, *http.Request)

func (invoke httpHandlerFuncInvoker) Invoke(args []interface{}) ([]reflect.Value, error) {
	invoke(args[0].(http.ResponseWriter), args[1].(*http.Request))
	return nil, nil
}

func fastInvokeWarpHandler(h HandlerFunc) HandlerFunc {
	if reflect.TypeOf(h).Kind() != reflect.Func {
		panic(fmt.Sprintf("handler must be a callable function, but got %T", h))
	}

	if inject.IsFastInvoker(h) {
		return h
	}

	switch v := h.(type) {
	case func(*Context):
		return ContextInvoker(v)
	case func(http.ResponseWriter, *http.Request):
		return httpHandlerFuncInvoker(v)
	case http.HandlerFunc:
		return httpHandlerFuncInvoker(v)
	}
	return h
}

func fastInvokeWarpHandlerChain(hc HandlersChain) {
	for i, handlerFunc := range hc {
		hc[i] = fastInvokeWarpHandler(handlerFunc)
	}
}
