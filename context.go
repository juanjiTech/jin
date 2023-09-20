package jin

import (
	"fmt"
	"github.com/juanjiTech/inject/v2"
	"math"
	"net/http"
)

// abortIndex represents a typical value used in abort functions.
const abortIndex int8 = math.MaxInt8 >> 1

type Context struct {
	inject.Injector
	writermem responseWriter
	Request   *http.Request
	Writer    ResponseWriter
	Params    Params

	handlers HandlersChain
	fullPath string
	index    int8

	params       *Params
	skippedNodes *[]skippedNode
}

func (c *Context) reset() {
	c.Injector.Reset()
	c.Writer = &c.writermem
	c.Params = c.Params[:0]

	*c.params = (*c.params)[:0]
	c.handlers = nil
	c.index = -1
	c.fullPath = ""
}

// FullPath returns a matched route full path. For not found routes
// returns an empty string.
//
//	router.GET("/user/:id", func(c *gin.Context) {
//	    c.FullPath() == "/user/:id" // true
//	})
func (c *Context) FullPath() string {
	return c.fullPath
}

// Next should be used only inside middleware.
// It executes the pending handlers in the chain inside the calling handler.
// See example in GitHub.
func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		h := c.handlers[c.index]
		values, err := c.Invoke(h)
		if err != nil {
			panic(fmt.Sprintf("unable to invoke the %s handler [%s:%T]: %v",
				ordinalize(int(c.index)), nameOfFunction(h), h, err))
		}
		c.index++

		for _, val := range values {
			c.Map(val.Interface())
		}
	}
}

// IsAborted returns true if the current context was aborted.
func (c *Context) IsAborted() bool {
	return c.index >= abortIndex
}

// Abort prevents pending handlers from being called. Note that this will not stop the current handler.
// Let's say you have an authorization middleware that validates that the current request is authorized.
// If the authorization fails (ex: the password does not match), call Abort to ensure the remaining handlers
// for this request are not called.
func (c *Context) Abort() {
	c.index = abortIndex
}
