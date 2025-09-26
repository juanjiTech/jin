package testcase

import (
	"fmt"
	"io"
	"testing"

	"github.com/juanjiTech/jin"
	"github.com/juanjiTech/jin/middleware/binding"
	"github.com/juanjiTech/jin/render"
)

type ExampleReq struct {
	Name string `json:"name"`
}

func TestJinBind(t *testing.T) {
	engine := jin.New()
	engine.POST("/json/1", binding.JSON(ExampleReq{}), func(ctx *jin.Context, req ExampleReq) {
		_ = render.WriteJSON(ctx.Writer, req.Name)
		fmt.Println("Received Name: ", req.Name)
		buf, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			fmt.Printf("Failed to read request body: %v\n", err)
		}
		fmt.Println("Request Body:", string(buf), "End")
	})
	_ = engine.Run(":8080")
}
