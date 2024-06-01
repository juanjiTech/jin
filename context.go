package jin

import (
	"fmt"
	"github.com/juanjiTech/inject/v2"
	"github.com/juanjiTech/jin/render"
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

	// Errors is a list of errors attached to all the handlers/middlewares who used this context.
	Errors []*error

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

	c.Errors = c.Errors[:0]

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
		if h == nil {
			c.index++
			continue
		}
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

// Error attaches an error to the current context. The error is pushed to a list of errors.
// It's a good idea to call Error for each error that occurred during the resolution of a request.
// A middleware can be used to collect all the errors and push them to a database together,
// print a log, or append it in the HTTP response.
// Error will panic if err is nil.
func (c *Context) Error(err error) *error {
	if err == nil {
		panic("err is nil")
	}

	c.Errors = append(c.Errors, &err)
	return &err
}

// Status sets the HTTP response code.
func (c *Context) Status(code int) {
	c.Writer.WriteHeader(code)
}

// bodyAllowedForStatus is a copy of http.bodyAllowedForStatus non-exported function.
func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == http.StatusNoContent:
		return false
	case status == http.StatusNotModified:
		return false
	}
	return true
}

// Render writes the response headers and calls render.Render to render data.
func (c *Context) Render(code int, r render.Render) {
	c.Status(code)

	if !bodyAllowedForStatus(code) {
		r.WriteContentType(c.Writer)
		c.Writer.WriteHeaderNow()
		return
	}

	if err := r.Render(c.Writer); err != nil {
		// Pushing error to c.Errors
		_ = c.Error(err)
		c.Abort()
	}
}
