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
	Query string `query:"query"`
	Value int    `query:"value"`
}

func TestJinQuery(t *testing.T) {
	engine := jin.New()
	engine.POST("/query/1", binding.Query(ExampleQuery{}), func(ctx *jin.Context, req ExampleQuery) {
		_ = render.WriteJSON(ctx.Writer, map[string]int{
			req.Query: req.Value,
		})
		fmt.Println("Received Query: ", req.Query)
		fmt.Println("Received Value: ", req.Value)
		buf, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			fmt.Printf("Failed to read request body: %v\n", err)
		}
		fmt.Println("Request Body:", string(buf), "End")
	})
	engine.GET("/query/2", binding.Query(ExampleQuery{}), func(ctx *jin.Context, req ExampleQuery) {
		_ = render.WriteJSON(ctx.Writer, map[string]int{
			req.Query: req.Value,
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
