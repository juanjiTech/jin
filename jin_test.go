package jin

import (
	"github.com/juanjiTech/inject/v2"
	"net/http"
	"strconv"
	"testing"
)

func TestDefaultEngine(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	SetMode(DebugMode)
	i := inject.New()
	i.Map(123)
	e := Default()
	e.SetParent(i)
	e.NoRoute(http.NotFound)
	e.GET("/ping", func(c *Context) {
		_, _ = c.Writer.WriteString("pong")
	}, func(c *Context, v int) {
		v++
		_, _ = c.Writer.WriteString(strconv.Itoa(v))
	})
	e.GET("/panic", func() {
		panic("I am panic")
	})

	_ = e.Run(":8080")
}
