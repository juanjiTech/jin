package jin

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// CreateTestContext returns a fresh engine and context for testing purposes
func CreateTestContext(w http.ResponseWriter) (c *Context, r *Engine) {
	r = New()
	c = r.allocateContext(0)
	c.reset()
	c.writermem.reset(w)
	return
}

// CreateTestContextOnly returns a fresh context base on the engine for testing purposes
func CreateTestContextOnly(w http.ResponseWriter, r *Engine) (c *Context) {
	c = r.allocateContext(r.maxParams)
	c.reset()
	c.writermem.reset(w)
	return
}

func compareFunc(t *testing.T, a, b any) {
	sf1 := reflect.ValueOf(a)
	sf2 := reflect.ValueOf(b)
	if sf1.Pointer() != sf2.Pointer() {
		t.Error("different functions")
	}
}

func TestContextHandlers(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	assert.Nil(t, c.handlers)
	assert.Nil(t, c.handlers.Last())

	c.handlers = HandlersChain{}
	assert.NotNil(t, c.handlers)
	assert.Nil(t, c.handlers.Last())

	f := func(c *Context) {}
	g := func(c *Context) {}

	c.handlers = HandlersChain{f}
	compareFunc(t, f, c.handlers.Last())

	c.handlers = HandlersChain{f, g}
	compareFunc(t, g, c.handlers.Last())
}

func TestContextAbort(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	assert.False(t, c.IsAborted())

	c.Abort()
	assert.True(t, c.IsAborted())
	assert.Equal(t, abortIndex, c.index)
}

func TestContextError(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	assert.Empty(t, c.Errors)

	err := errors.New("test error")
	c.Error(err)
	assert.Len(t, c.Errors, 1)
	assert.Equal(t, "test error", (*c.Errors[0]).Error())

	assert.Panics(t, func() {
		c.Error(nil)
	})
}

func TestContextStatus(t *testing.T) {
	recorder := httptest.NewRecorder()
	c, _ := CreateTestContext(recorder)
	c.Status(http.StatusTeapot)
	c.Writer.WriteHeaderNow()
	assert.Equal(t, http.StatusTeapot, recorder.Code)
}

func TestContextFullPath(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.fullPath = "/test/path"
	assert.Equal(t, "/test/path", c.FullPath())
}
