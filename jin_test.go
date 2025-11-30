package jin

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/juanjiTech/inject/v2"
	"github.com/stretchr/testify/assert"
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
		return "", errors.New("damn this is an error!")
	}, respHandler)
	e.GET("/noerror", func() (string, error) {
		return "no error!", nil
	}, respHandler)

	e.GET("/panic", func() {
		panic("I am panic")
	})

	_ = e.Run(":8080")
}

func TestDefault(t *testing.T) {
	engine := Default()
	assert.NotNil(t, engine)
	assert.Equal(t, 1, len(engine.Handlers))
}

func TestEngineHandler(t *testing.T) {
	// Test with UseH2C = false
	engine := New()
	handler := engine.Handler()
	assert.Equal(t, engine, handler)

	// Test with UseH2C = true
	engine.UseH2C = true
	handler = engine.Handler()
	assert.NotEqual(t, engine, handler)
}

func TestEngineNoMethod(t *testing.T) {
	engine := New()
	engine.NoMethod(func(c *Context) {
		c.Status(http.StatusTeapot)
	})
	assert.Equal(t, 1, len(engine.noMethod))
}

func TestEngineRedirectTrailingSlash(t *testing.T) {
	engine := New()
	engine.RedirectTrailingSlash = true
	engine.GET("/foo", func(c *Context) {
		c.Writer.Write([]byte("foo"))
	})
	engine.GET("/bar/", func(c *Context) {
		c.Writer.Write([]byte("bar"))
	})
	engine.POST("/bar/", func(c *Context) {
		c.Writer.Write([]byte("bar"))
	})

	t.Run("path with trailing slash", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/foo/", nil)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusMovedPermanently, w.Code)
		assert.Equal(t, "/foo", w.Header().Get("Location"))
	})

	t.Run("path without trailing slash", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/bar", nil)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusMovedPermanently, w.Code)
		assert.Equal(t, "/bar/", w.Header().Get("Location"))
	})

	t.Run("POST path without trailing slash", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/bar", nil)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		assert.Equal(t, "/bar/", w.Header().Get("Location"))
	})

	t.Run("path with trailing slash and X-Forwarded-Prefix", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/foo/", nil)
		req.Header.Set("X-Forwarded-Prefix", "/prefix")
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusMovedPermanently, w.Code)
		assert.Equal(t, "/prefix/foo", w.Header().Get("Location"))
	})
}

func TestEngineRedirectFixedPath(t *testing.T) {
	engine := New()
	engine.RedirectFixedPath = true
	engine.GET("/foo", func(c *Context) {
		c.Writer.Write([]byte("foo"))
	})
	engine.POST("/bar", func(c *Context) {
		c.Writer.Write([]byte("bar"))
	})

	t.Run("case-insensitive redirect", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/FOO", nil)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusMovedPermanently, w.Code)
		assert.Equal(t, "/foo", w.Header().Get("Location"))
	})

	t.Run("path cleaning redirect", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/..//foo", nil)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusMovedPermanently, w.Code)
		assert.Equal(t, "/foo", w.Header().Get("Location"))
	})

	t.Run("POST request redirect", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/BAR", nil)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		assert.Equal(t, "/bar", w.Header().Get("Location"))
	})

	t.Run("no redirect for non-existent path", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/nonexistent", nil)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestEngineRun(t *testing.T) {
	engine := New()
	engine.GET("/", func(c *Context) {
		c.Writer.Write([]byte("ok"))
	})

	server := &http.Server{
		Addr:    ":8089",
		Handler: engine,
	}

	go func() {
		// The server always returns a non-nil error.
		_ = server.ListenAndServe()
	}()

	// Wait for the server to start
	time.Sleep(5 * time.Millisecond)

	resp, err := http.Get("http://localhost:8089/")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "ok", string(body))

	// Shutdown the server
	err = server.Shutdown(context.Background())
	assert.NoError(t, err)
}

func TestEngineRunError(t *testing.T) {
	// Create a listener to occupy a port
	listener, err := net.Listen("tcp", ":0")
	assert.NoError(t, err)
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	addr := ":" + strconv.Itoa(port)

	engine := New()
	errCh := make(chan error)
	go func() {
		errCh <- engine.Run(addr)
	}()

	// Wait for the Run function to return an error
	err = <-errCh
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "address already in use")
}
