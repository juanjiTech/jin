package testcase

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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
