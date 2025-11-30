package testcase

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/juanjiTech/jin"
	"github.com/juanjiTech/jin/middleware/binding"
	"github.com/juanjiTech/jin/render"
)

type ExampleQuery struct {
	Query           string     `query:"query"`
	Value           int        `query:"value"`
	TestIntSlice    []int      `query:"test_int_slice"`
	TestStringSlice []string   `query:"test_string_slice"`
	TestStruct      TestStruct `query:"teststruct"`
}

type TestStruct struct {
	Field1 int    `query:"field1"`
	Field2 string `query:"field2"`
}

func TestJinQuery(t *testing.T) {
	engine := jin.New()
	engine.POST("/query/1", binding.Query(ExampleQuery{}), func(ctx *jin.Context, req ExampleQuery) {
		_ = render.WriteJSON(ctx.Writer, map[string]any{
			"requests": req,
		})
	})
	engine.GET("/query/2", binding.Query(&ExampleQuery{}), func(ctx *jin.Context, req ExampleQuery) {
		_ = render.WriteJSON(ctx.Writer, map[string]any{
			"requests": req,
		})
	})

	t.Run("POST request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/query/1?query=test&value=123&test_int_slice=1,2,3&test_string_slice=a,b,c&teststruct.field1=456&teststruct.field2=def", nil)
		engine.ServeHTTP(w, req)

		if http.StatusOK != w.Code {
			t.Errorf("expected status code %d, but got %d", http.StatusOK, w.Code)
		}

		var resp map[string]ExampleQuery
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response body: %v", err)
		}

		expected := ExampleQuery{
			Query:           "test",
			Value:           123,
			TestIntSlice:    []int{1, 2, 3},
			TestStringSlice: []string{"a", "b", "c"},
			TestStruct: TestStruct{
				Field1: 456,
				Field2: "def",
			},
		}

		if !reflect.DeepEqual(expected, resp["requests"]) {
			t.Errorf("expected %+v, but got %+v", expected, resp["requests"])
		}
	})

	t.Run("GET request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/query/2?query=test2&value=456", nil)
		engine.ServeHTTP(w, req)

		if http.StatusOK != w.Code {
			t.Errorf("expected status code %d, but got %d", http.StatusOK, w.Code)
		}

		var resp map[string]ExampleQuery
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response body: %v", err)
		}

		expected := ExampleQuery{
			Query: "test2",
			Value: 456,
		}

		// We only check the fields that were sent, the others should be zero-valued
		actual := resp["requests"]
		if expected.Query != actual.Query || expected.Value != actual.Value {
			t.Errorf("expected query %s and value %d, but got query %s and value %d", expected.Query, expected.Value, actual.Query, actual.Value)
		}
	})
}
