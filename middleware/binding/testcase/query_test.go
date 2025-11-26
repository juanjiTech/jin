package testcase

import (
	"fmt"
	"io"
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
			req.Query:  req.Value,
			"requests": req,
		})
		fmt.Println("Received Query: ", req.Query)
		fmt.Println("Received Value: ", req.Value)
		buf, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			fmt.Printf("Failed to read request body: %v\n", err)
		}
		fmt.Println("Request Body:", string(buf), "End")
	})
	engine.GET("/query/2", binding.Query(&ExampleQuery{}), func(ctx *jin.Context, req ExampleQuery) {
		_ = render.WriteJSON(ctx.Writer, map[string]any{
			req.Query:  req.Value,
			"requests": req,
		})
		fmt.Println("Received Query: ", req.Query)
		fmt.Println("Received Value: ", req.Value)
		buf, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			fmt.Printf("Failed to read request body: %v\n", err)
		}
		fmt.Println("Request Body:", string(buf), "End")
	})
	_ = engine.Run(":8080")
}
