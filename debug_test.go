package jin

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDebugging(t *testing.T) {
	SetMode(DebugMode)
	assert.True(t, IsDebugging())
	SetMode(ReleaseMode)
	assert.False(t, IsDebugging())
	SetMode(TestMode)
	assert.False(t, IsDebugging())
}

func TestDebugPrint(t *testing.T) {
	oldWriter := DefaultWriter
	r, w, _ := os.Pipe()
	DefaultWriter = w
	SetMode(DebugMode)

	debugPrint("hello %s", "world")
	w.Close()
	outputBytes, _ := io.ReadAll(r)
	output := string(outputBytes)
	assert.Contains(t, output, "[JIN-debug] hello world")

	r, w, _ = os.Pipe()
	DefaultWriter = w
	debugPrintWARNINGNew()
	w.Close()
	outputBytes, _ = io.ReadAll(r)
	output = string(outputBytes)
	assert.Contains(t, output, "[JIN-debug] [WARNING] Running in \"debug\" mode.")

	DefaultWriter = oldWriter
	SetMode(TestMode)
}

func TestDebugPrintError(t *testing.T) {
	oldErrorWriter := DefaultErrorWriter
	r, w, _ := os.Pipe()
	DefaultErrorWriter = w
	SetMode(DebugMode)

	err := errors.New("test error")
	debugPrintError(err)
	w.Close()
	outputBytes, _ := io.ReadAll(r)
	output := string(outputBytes)
	assert.Contains(t, output, "[JIN-debug] [ERROR] test error")

	r, w, _ = os.Pipe()
	DefaultErrorWriter = w
	debugPrintError(nil)
	w.Close()
	outputBytes, _ = io.ReadAll(r)
	output = string(outputBytes)
	assert.Equal(t, "", output)

	DefaultErrorWriter = oldErrorWriter
	SetMode(TestMode)
}

func TestDebugPrintRoute(t *testing.T) {
	oldWriter := DefaultWriter
	var out bytes.Buffer
	DefaultWriter = &out
	SetMode(DebugMode)

	handler := func(c *Context) {}
	handlers := HandlersChain{handler}
	debugPrintRoute("GET", "/test", handlers)
	assert.Contains(t, out.String(), "[JIN-debug] GET")
	assert.Contains(t, out.String(), "/test")
	assert.Contains(t, out.String(), "TestDebugPrintRoute.func1")
	assert.Contains(t, out.String(), "(1 handlers)")
	out.Reset()

	DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		DefaultWriter.Write([]byte(strings.Join([]string{httpMethod, absolutePath, handlerName, "custom"}, " ")))
	}
	debugPrintRoute("POST", "/custom", handlers)
	assert.Contains(t, out.String(), "POST /custom")
	assert.Contains(t, out.String(), "TestDebugPrintRoute.func1")
	assert.Contains(t, out.String(), "custom")

	DefaultWriter = oldWriter
	DebugPrintRouteFunc = nil
	SetMode(TestMode)
}
