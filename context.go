package jin

import (
	"fmt"
	"github.com/juanjiTech/inject"
	"net/http"
)

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
