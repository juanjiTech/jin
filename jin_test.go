package jin

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/juanjiTech/inject/v2"
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

	respHandler := func(c *Context, v string, err error) {
		t.Log("response handler:", v)
		if err != nil {
			t.Log("response error:", err)
			_, _ = c.Writer.WriteString("error: " + err.Error())
			return
		}
		_, _ = c.Writer.WriteString(v)
		return
	}

	e.GET("/pong", func(v int) string {
		t.Log("pong", v)
		return "pong " + strconv.Itoa(v)
	}, func(c string) {
		t.Log("final", c)
	}, respHandler)
	e.GET("/error", func() (string, error) {
		return "", fmt.Errorf("damn this is an error!")
	}, respHandler)
	e.GET("/noerror", func() (string, error) {
		return "no error!", nil
	}, respHandler)

	e.GET("/panic", func() {
		panic("I am panic")
	})

	_ = e.Run(":8080")
}
