package jin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouterGroup(t *testing.T) {
	engine := New()
	api := engine.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/hello", func(c *Context) {
		c.Writer.WriteString("hello")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/hello", nil)
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, but got %d", http.StatusOK, w.Code)
	}
	if w.Body.String() != "hello" {
		t.Errorf("expected body %q, but got %q", "hello", w.Body.String())
	}
}

func TestRouterMiddleware(t *testing.T) {
	engine := New()
	group := engine.Group("/api")

	// Middleware 1: adds a header
	group.Use(func(c *Context) {
		c.Writer.Header().Set("X-Test-Middleware1", "true")
		c.Next()
	})

	// Middleware 2: writes to the body
	group.Use(func(c *Context) {
		c.Writer.WriteString("middleware2-start|")
		c.Next()
		c.Writer.WriteString("|middleware2-end")
	})

	group.GET("/hello", func(c *Context) {
		c.Writer.WriteString("hello")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/hello", nil)
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, but got %d", http.StatusOK, w.Code)
	}
	if w.Header().Get("X-Test-Middleware1") != "true" {
		t.Errorf("expected header X-Test-Middleware1 to be 'true'")
	}
	expectedBody := "middleware2-start|hello|middleware2-end"
	if w.Body.String() != expectedBody {
		t.Errorf("expected body %q, but got %q", expectedBody, w.Body.String())
	}
}

func TestRouterHTTPMethods(t *testing.T) {
	engine := New()

	handler := func(c *Context) {
		c.Writer.WriteString(c.Request.Method)
	}

	engine.POST("/test", handler)
	engine.PUT("/test", handler)
	engine.DELETE("/test", handler)
	engine.PATCH("/test", handler)
	engine.OPTIONS("/test", handler)
	engine.HEAD("/test", handler)

	methods := []string{"POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"}

	for _, method := range methods {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(method, "/test", nil)
		engine.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("for method %s, expected status code %d, but got %d", method, http.StatusOK, w.Code)
		}
		if method != "HEAD" {
			if w.Body.String() != method {
				t.Errorf("for method %s, expected body %q, but got %q", method, method, w.Body.String())
			}
		}
	}
}

func TestRouterAnyMethod(t *testing.T) {
	engine := New()
	engine.Any("/any", func(c *Context) {
		c.Writer.WriteString("any")
	})

	methods := []string{
		http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch,
		http.MethodHead, http.MethodOptions, http.MethodDelete,
	}

	for _, method := range methods {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(method, "/any", nil)
		engine.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("for method %s, expected status code %d, but got %d", method, http.StatusOK, w.Code)
		}
		if method != "HEAD" {
			if w.Body.String() != "any" {
				t.Errorf("for method %s, expected body 'any', but got %q", method, w.Body.String())
			}
		}
	}
}

func TestRecovery(t *testing.T) {
	engine := New()
	engine.Use(Recovery())
	engine.GET("/panic", func(c *Context) {
		panic("test panic")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/panic", nil)
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status code %d, but got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestNoRoute(t *testing.T) {
	t.Run("404 Not Found", func(t *testing.T) {
		engine := New()
		engine.NoRoute(func(c *Context) {
			c.Writer.WriteHeader(http.StatusNotFound)
			c.Writer.WriteString("custom 404")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/non-existent", nil)
		engine.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status code %d, but got %d", http.StatusNotFound, w.Code)
		}
		if w.Body.String() != "custom 404" {
			t.Errorf("expected body 'custom 404', but got %q", w.Body.String())
		}
	})

	t.Run("405 Method Not Allowed", func(t *testing.T) {
		engine := New()
		engine.HandleMethodNotAllowed = true
		engine.POST("/exists", func(c *Context) {})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/exists", nil)
		engine.ServeHTTP(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status code %d, but got %d", http.StatusMethodNotAllowed, w.Code)
		}
	})
}

func TestRouterGroupBasePath(t *testing.T) {
	engine := New()
	assert.Equal(t, "/", engine.BasePath())

	api := engine.Group("/api")
	assert.Equal(t, "/api", api.BasePath())

	v1 := api.Group("/v1")
	assert.Equal(t, "/api/v1", v1.BasePath())
}

func TestRouterHandleAndMatch(t *testing.T) {
	handler := func(c *Context) {
		c.Writer.WriteString(c.Request.Method)
	}

	t.Run("TestHandle", func(t *testing.T) {
		engine := New()
		engine.Handle("CUSTOM", "/custom", handler)
		engine.Handle("CUSTOM", "/custom2", handler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("CUSTOM", "/custom", nil)
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "CUSTOM", w.Body.String())

		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest("CUSTOM", "/custom2", nil)
		engine.ServeHTTP(w2, req2)

		assert.Equal(t, http.StatusOK, w2.Code)
		assert.Equal(t, "CUSTOM", w2.Body.String())
	})

	t.Run("TestHandleInvalidMethod", func(t *testing.T) {
		engine := New()
		handler := func(c *Context) {}
		assert.Panics(t, func() {
			engine.Handle("invalid method", "/test", handler)
		})
		assert.Panics(t, func() {
			engine.Handle("INVALID-METHOD", "/test", handler)
		})
	})

	t.Run("TestMatch", func(t *testing.T) {
		engine := New()
		engine.Match([]string{"GET", "POST"}, "/match", handler)

		// Test GET
		wGet := httptest.NewRecorder()
		reqGet, _ := http.NewRequest("GET", "/match", nil)
		engine.ServeHTTP(wGet, reqGet)
		assert.Equal(t, http.StatusOK, wGet.Code)
		assert.Equal(t, "GET", wGet.Body.String())

		// Test POST
		wPost := httptest.NewRecorder()
		reqPost, _ := http.NewRequest("POST", "/match", nil)
		engine.ServeHTTP(wPost, reqPost)
		assert.Equal(t, http.StatusOK, wPost.Code)
		assert.Equal(t, "POST", wPost.Body.String())
	})
}
