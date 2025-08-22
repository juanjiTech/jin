package binding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/juanjiTech/jin"
	"io"
)

func noSideEffectJsonUnmarshal(reader io.ReadCloser, v any) (bytes.Buffer, error) {
	var buf bytes.Buffer
	teeReader := io.TeeReader(reader, &buf)
	decoder := json.NewDecoder(teeReader)
	decoder.UseNumber()
	if err := decoder.Decode(v); err != nil {
		return buf, fmt.Errorf("string: `%s`, error: `%w`", buf.String(), err)
	}
	return buf, reader.Close()
}

func JSON[T any](model T) jin.HandlerFunc {
	_ = model // to avoid unused variable warning
	return func(ctx *jin.Context) {
		var t T
		r := ctx.Request
		if r.Body != nil {
			buf, err := noSideEffectJsonUnmarshal(r.Body, &t)
			if err != nil {
				ctx.Error(err)
				ctx.Abort()
				return
			}
			r.Body = io.NopCloser(&buf)
			ctx.Map(t)
		}
		ctx.Next()
	}
}
