package binding

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/juanjiTech/jin"
	"github.com/juanjiTech/jin/render"
	"github.com/stretchr/testify/assert"
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
	engine.POST("/query/1", Query(ExampleQuery{}), func(ctx *jin.Context, req ExampleQuery) {
		_ = render.WriteJSON(ctx.Writer, map[string]any{
			"requests": req,
		})
	})
	engine.GET("/query/2", Query(&ExampleQuery{}), func(ctx *jin.Context, req ExampleQuery) {
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

func TestQueryBindingPanic(t *testing.T) {
	assert.Panics(t, func() {
		// Pass a non-struct type to Query
		Query(123)
	}, "Query should panic if model is not a struct")
}

type queryAnonymous struct {
	// This field will be part of the anonymous struct
	AnonymousField string `query:"anonymous_field"`
}

type queryUnsupportedFields struct {
	// Anonymous field, should be skipped
	queryAnonymous
	// No 'query' tag, should be skipped
	NoTag string
	// Unsupported type, should be skipped
	Unsupported chan int `query:"unsupported"`
	// Supported field for control
	Foo string `query:"foo"`
}

func TestQueryBindingWithUnsupportedFields(t *testing.T) {
	e := jin.New()
	e.GET("/", Query(queryUnsupportedFields{}), func(c *jin.Context, model queryUnsupportedFields) {
		// Check that only the supported field was set
		assert.Equal(t, "bar", model.Foo)
		// Check that skipped fields are empty/zero because anonymous fields are skipped
		assert.Empty(t, model.AnonymousField)
		assert.Empty(t, model.NoTag)
		assert.Nil(t, model.Unsupported)

		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/?foo=bar&unsupported=baz&anonymous_field=anon", nil)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

type comprehensiveQuery struct {
	BoolValue  bool      `query:"bool_value"`
	UintValue  uint      `query:"uint_value"`
	FloatValue float64   `query:"float_value"`
	BoolSlice  []bool    `query:"bool_slice"`
	UintSlice  []uint    `query:"uint_slice"`
	FloatSlice []float64 `query:"float_slice"`
}

func TestQueryBindingComprehensiveTypes(t *testing.T) {
	engine := jin.New()
	engine.GET("/comprehensive", Query(comprehensiveQuery{}), func(ctx *jin.Context, req comprehensiveQuery) {
		_ = render.WriteJSON(ctx.Writer, req)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/comprehensive?bool_value=true&uint_value=123&float_value=45.67&bool_slice=true,false,1,0&uint_slice=1,2,3&float_slice=1.1,2.2,3.3", nil)
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp comprehensiveQuery
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	expected := comprehensiveQuery{
		BoolValue:  true,
		UintValue:  123,
		FloatValue: 45.67,
		BoolSlice:  []bool{true, false, true, false},
		UintSlice:  []uint{1, 2, 3},
		FloatSlice: []float64{1.1, 2.2, 3.3},
	}

	assert.Equal(t, expected, resp)
}
