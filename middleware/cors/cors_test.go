package cors

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/juanjiTech/jin"
	"github.com/stretchr/testify/assert"
)

// TestCorsDefault tests the default CORS middleware which allows all origins.
func TestCorsDefault(t *testing.T) {
	e := jin.New()
	e.Use(Default())
	e.GET("/", func(c *jin.Context) {
		c.Writer.WriteString("ok")
	})
	e.OPTIONS("/", func(c *jin.Context) {}) // Handler for OPTIONS is required

	// Test Preflight request
	req := httptest.NewRequest("OPTIONS", "/", nil)
	req.Host = "api.example.com"
	req.Header.Set("Origin", "http://example.com")
	req.Header.Set("Access-Control-Request-Method", "GET")
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Methods"))

	// Test Normal GET request
	req = httptest.NewRequest("GET", "/", nil)
	req.Host = "api.example.com"
	req.Header.Set("Origin", "http://example.com")
	w = httptest.NewRecorder()
	e.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
}

// TestCorsSpecificOrigins tests CORS with specific origins allowed.
func TestCorsSpecificOrigins(t *testing.T) {
	config := DefaultConfig()
	config.AllowOrigins = []string{"http://example.com", "http://*.example.org"}
	config.AllowWildcard = true

	e := jin.New()
	e.Use(New(config))
	e.GET("/", func(c *jin.Context) {
		c.Writer.WriteString("ok")
	})
	e.OPTIONS("/", func(c *jin.Context) {})

	// Test allowed origin
	req := httptest.NewRequest("GET", "/", nil)
	req.Host = "api.other.com"
	req.Header.Set("Origin", "http://example.com")
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "http://example.com", w.Header().Get("Access-Control-Allow-Origin"))

	// Test allowed wildcard origin
	req = httptest.NewRequest("GET", "/", nil)
	req.Host = "api.other.com"
	req.Header.Set("Origin", "http://sub.example.org")
	w = httptest.NewRecorder()
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "http://sub.example.org", w.Header().Get("Access-Control-Allow-Origin"))

	// Test disallowed origin
	req = httptest.NewRequest("GET", "/", nil)
	req.Host = "api.other.com"
	req.Header.Set("Origin", "http://disallowed.com")
	w = httptest.NewRecorder()
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

// TestCorsSameOrigin tests that CORS headers are not added for same-origin requests.
func TestCorsSameOrigin(t *testing.T) {
	e := jin.New()
	e.Use(Default())
	e.GET("/", func(c *jin.Context) {
		c.Writer.WriteString("ok")
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Host = "example.com"
	req.Header.Set("Origin", "http://example.com")
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

// TestCorsWithOriginFunc tests the AllowOriginFunc configuration.
func TestCorsWithOriginFunc(t *testing.T) {
	config := DefaultConfig()
	config.AllowOriginFunc = func(origin string) bool {
		return origin == "http://allowed.com"
	}

	e := jin.New()
	e.Use(New(config))
	e.GET("/", func(c *jin.Context) {
		c.Writer.WriteString("ok")
	})

	// Test allowed origin
	req := httptest.NewRequest("GET", "/", nil)
	req.Host = "api.other.com"
	req.Header.Set("Origin", "http://allowed.com")
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "http://allowed.com", w.Header().Get("Access-Control-Allow-Origin"))

	// Test disallowed origin
	req = httptest.NewRequest("GET", "/", nil)
	req.Host = "api.other.com"
	req.Header.Set("Origin", "http://disallowed.com")
	w = httptest.NewRecorder()
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

// TestConfigValidation tests the Validate function for the config.
func TestConfigValidation(t *testing.T) {
	// Conflict: AllowAllOrigins and AllowOrigins
	config := DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowOrigins = []string{"http://example.com"}
	err := config.Validate()
	assert.Error(t, err)

	// Conflict: All origins disabled
	config = Config{}
	err = config.Validate()
	assert.Error(t, err)

	// Bad origin schema
	config = DefaultConfig()
	config.AllowOrigins = []string{"example.com"}
	err = config.Validate()
	assert.Error(t, err)

	// Valid config
	config = DefaultConfig()
	config.AllowOrigins = []string{"http://example.com"}
	err = config.Validate()
	assert.NoError(t, err)
}

// TestCorsWithCredentials tests the AllowCredentials setting.
func TestCorsWithCredentials(t *testing.T) {
	config := DefaultConfig()
	config.AllowOrigins = []string{"http://example.com"}
	config.AllowCredentials = true

	e := jin.New()
	e.Use(New(config))
	e.GET("/", func(c *jin.Context) {
		c.Writer.WriteString("ok")
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Host = "api.other.com"
	req.Header.Set("Origin", "http://example.com")
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "http://example.com", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
}

func TestCorsAdvancedConfig(t *testing.T) {
	// Test AddAllowMethods, AddAllowHeaders, AddExposeHeaders
	config := DefaultConfig()
	config.AddAllowMethods("TRACE")
	config.AddAllowHeaders("X-Custom-Header")
	config.AddExposeHeaders("X-Custom-Response")
	config.AllowOrigins = []string{"http://example.com"}

	// Test different schemas
	config.AllowBrowserExtensions = true
	config.AllowWebSockets = true
	config.AllowFiles = true
	assert.True(t, config.validateAllowedSchemas("chrome-extension://test"))
	assert.True(t, config.validateAllowedSchemas("ws://test"))
	assert.True(t, config.validateAllowedSchemas("file://test"))

	// Test wildcard parsing and validation
	config.AllowWildcard = true
	config.AllowOrigins = []string{"http://*.example.com", "http://example.*", "http://prefix.*.suffix"}
	wRules := config.parseWildcardRules()
	assert.Len(t, wRules, 3)

	corsMiddleware := newCors(config)
	assert.True(t, corsMiddleware.validateWildcardOrigin("http://sub.example.com"))
	assert.True(t, corsMiddleware.validateWildcardOrigin("http://example.org"))
	assert.False(t, corsMiddleware.validateWildcardOrigin("http://another.domain"))

	// Test generateNormalHeaders with ExposeHeaders
	e := jin.New()
	e.Use(New(config))
	e.GET("/", func(c *jin.Context) {
		c.Writer.WriteString("ok")
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Host = "api.other.com"
	req.Header.Set("Origin", "http://sub.example.com")
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Access-Control-Expose-Headers"), "X-Custom-Response")
}

// TestConfigPanic tests panic conditions when creating a new middleware.
func TestConfigPanic(t *testing.T) {
	// Panic on invalid config
	assert.Panics(t, func() {
		config := Config{} // All origins disabled
		New(config)
	})

	// Panic on multiple wildcards
	assert.Panics(t, func() {
		config := DefaultConfig()
		config.AllowWildcard = true
		config.AllowOrigins = []string{"http://*.example.*"}
		New(config)
	})
}
