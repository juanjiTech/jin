package render

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Foo string `json:"foo" yaml:"foo"`
	Bar int    `json:"bar" yaml:"bar"`
}

func TestJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := testStruct{Foo: "bar", Bar: 1}
	err := JSON{Data: data}.Render(w)

	assert.NoError(t, err)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	var result testStruct
	err = json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, data, result)
}

func TestIndentedJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := testStruct{Foo: "bar", Bar: 1}
	err := IndentedJSON{Data: data}.Render(w)

	assert.NoError(t, err)
	expectedBody := "{\n    \"foo\": \"bar\",\n    \"bar\": 1\n}"
	assert.Equal(t, expectedBody, w.Body.String())
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	// Test WriteContentType
	w = httptest.NewRecorder()
	JSON{}.WriteContentType(w)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestSecureJSON(t *testing.T) {
	// Test with array data
	w := httptest.NewRecorder()
	dataArray := []testStruct{{Foo: "bar", Bar: 1}}
	err := SecureJSON{Prefix: "while(1);", Data: dataArray}.Render(w)
	assert.NoError(t, err)
	expectedBodyArray := "while(1);[{\"foo\":\"bar\",\"bar\":1}]"
	assert.Equal(t, expectedBodyArray, w.Body.String())
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	// Test with object data (should not add prefix)
	w = httptest.NewRecorder()
	dataObject := testStruct{Foo: "bar", Bar: 1}
	err = SecureJSON{Prefix: "while(1);", Data: dataObject}.Render(w)
	assert.NoError(t, err)
	expectedBodyObject := "{\"foo\":\"bar\",\"bar\":1}"
	assert.Equal(t, expectedBodyObject, w.Body.String())
}

func TestJsonpJSON(t *testing.T) {
	// Test with callback
	w := httptest.NewRecorder()
	data := testStruct{Foo: "bar", Bar: 1}
	err := JsonpJSON{Callback: "my_callback", Data: data}.Render(w)
	assert.NoError(t, err)
	expectedBody := "my_callback({\"foo\":\"bar\",\"bar\":1});"
	assert.Equal(t, expectedBody, w.Body.String())
	assert.Equal(t, "application/javascript; charset=utf-8", w.Header().Get("Content-Type"))

	// Test without callback
	w = httptest.NewRecorder()
	err = JsonpJSON{Data: data}.Render(w)
	assert.NoError(t, err)
	expectedBody = "{\"foo\":\"bar\",\"bar\":1}"
	assert.Equal(t, expectedBody, w.Body.String())
	// Content-Type should be application/javascript even without callback
	assert.Equal(t, "application/javascript; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestAsciiJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"lang": "GO语言"}
	err := AsciiJSON{Data: data}.Render(w)

	assert.NoError(t, err)
	expectedBody := "{\"lang\":\"GO\\u8bed\\u8a00\"}"
	assert.Equal(t, expectedBody, w.Body.String())
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

func TestPureJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"html": "<b>Hello</b>"}
	err := PureJSON{Data: data}.Render(w)

	assert.NoError(t, err)
	// Note: The default json.Encoder escapes HTML, PureJSON should not.
	expectedBody := "{\"html\":\"<b>Hello</b>\"}\n"
	assert.Equal(t, expectedBody, w.Body.String())
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestYAML(t *testing.T) {
	w := httptest.NewRecorder()
	data := testStruct{Foo: "bar", Bar: 1}
	err := YAML{Data: data}.Render(w)

	assert.NoError(t, err)
	assert.Equal(t, "application/yaml; charset=utf-8", w.Header().Get("Content-Type"))

	var result testStruct
	err = yaml.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, data, result)
}

func TestWriteContentType(t *testing.T) {
	w := httptest.NewRecorder()
	// Set a content type first
	w.Header().Set("Content-Type", "existing/type")
	// This should not overwrite the existing one
	writeContentType(w, []string{"application/json"})
	assert.Equal(t, "existing/type", w.Header().Get("Content-Type"))
}

// TestRenderWithError tests render methods with data that cannot be marshaled.
func TestRenderWithError(t *testing.T) {
	// Use a channel, which cannot be marshaled to JSON.
	data := make(chan int)

	wJson := httptest.NewRecorder()
	err := JSON{Data: data}.Render(wJson)
	assert.Error(t, err)

	wIndented := httptest.NewRecorder()
	err = IndentedJSON{Data: data}.Render(wIndented)
	assert.Error(t, err)

	wSecure := httptest.NewRecorder()
	err = SecureJSON{Data: data}.Render(wSecure)
	assert.Error(t, err)

	wJsonp := httptest.NewRecorder()
	err = JsonpJSON{Data: data}.Render(wJsonp)
	assert.Error(t, err)

	wAscii := httptest.NewRecorder()
	err = AsciiJSON{Data: data}.Render(wAscii)
	assert.Error(t, err)

	wPure := httptest.NewRecorder()
	err = PureJSON{Data: data}.Render(wPure)
	assert.Error(t, err)

	// For YAML, the underlying library now returns an error instead of panicking.
	type unmarshalable struct {
		F func()
	}
	wYaml := httptest.NewRecorder()
	err = YAML{Data: unmarshalable{F: func() {}}}.Render(wYaml)
	assert.Error(t, err)
}
