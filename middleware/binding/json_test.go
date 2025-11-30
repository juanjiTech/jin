package binding

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/juanjiTech/jin"
	"github.com/juanjiTech/jin/render"
	"github.com/stretchr/testify/assert"
)

type ExampleReq struct {
	Name string `json:"name"`
}

func TestJinBind(t *testing.T) {
	engine := jin.New()
	engine.POST("/json/1", JSON(ExampleReq{}), func(ctx *jin.Context, req ExampleReq) {
		_ = render.WriteJSON(ctx.Writer, req.Name)
		fmt.Println("Received Name: ", req.Name)
		buf, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			fmt.Printf("Failed to read request body: %v\n", err)
		}
		fmt.Println("Request Body:", string(buf), "End")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/json/1", strings.NewReader(`{"name": "test"}`))
	req.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(w, req)

	if http.StatusOK != w.Code {
		t.Errorf("expected status code %d, but got %d", http.StatusOK, w.Code)
	}
	expectedBody := `"test"`
	if expectedBody != w.Body.String() {
		t.Errorf("expected body %s, but got %s", expectedBody, w.Body.String())
	}
}

func TestJSONBindingWithInvalidJSON(t *testing.T) {
	e := jin.New()
	handlerCalled := false
	e.POST("/", JSON(ExampleReq{}), func(c *jin.Context) {
		handlerCalled = true
	})

	req := httptest.NewRequest("POST", "/", strings.NewReader(`{"name": "test",`)) // Invalid JSON
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)

	// The middleware should abort because of the JSON error, so the handler is not called.
	assert.False(t, handlerCalled)
}

func TestJSONBindingWithNilBody(t *testing.T) {
	e := jin.New()
	handlerCalled := false
	e.POST("/", JSON(ExampleReq{}), func(c *jin.Context) {
		handlerCalled = true
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("POST", "/", nil) // Nil body
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)

	// The middleware should just call Next(), so the handler is called.
	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, handlerCalled)
}

type errorReader struct{}

func (r *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func TestNoSideEffectJsonUnmarshalError(t *testing.T) {
	reader := &errorReader{}
	var v any
	_, err := noSideEffectJsonUnmarshal(io.NopCloser(reader), &v)
	assert.Error(t, err)
}
