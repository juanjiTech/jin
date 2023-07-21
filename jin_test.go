package jin

import (
	"github.com/juanjiTech/inject"
	"strconv"
	"testing"
)

func TestDefaultEngine(t *testing.T) {
	i := inject.New()
	i.Map(123)
	e := Default()
	e.SetParent(i)
	e.GET("/ping", func(c *Context) {
		_, _ = c.Writer.WriteString("pong")
	}, func(c *Context, v int) {
		v++
		_, _ = c.Writer.WriteString(strconv.Itoa(v))
	})

	e.Run(":8080")
}
