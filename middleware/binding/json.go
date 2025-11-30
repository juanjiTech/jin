package binding

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/juanjiTech/jin"
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
				// EOF means empty body, which is not an error in this context.
				// We should continue to the next handler.
				if !errors.Is(err, io.EOF) {
					ctx.Error(err)
					ctx.Abort()
					return
				}
			} else {
				// Only map the model if decoding was successful.
				ctx.Map(t)
			}
			r.Body = io.NopCloser(&buf)
		}
		ctx.Next()
	}
}
