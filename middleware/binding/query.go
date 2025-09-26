package binding

import (
	"net/url"
	"reflect"
	"strconv"

	"github.com/juanjiTech/jin"
)

func Query[T any](model T) jin.HandlerFunc {
	_ = model // to avoid unused variable warning
	typer := reflect.TypeOf(model)
	if typer.Kind() == reflect.Ptr {
		typer = typer.Elem()
	}
	// only support struct
	if typer.Kind() != reflect.Struct {
		panic("model must be a struct")
	}
	return func(ctx *jin.Context) {
		valueT := reflect.New(typer).Elem()
		u, _ := url.Parse(ctx.Request.RequestURI)
		queries := u.Query()
		for i := 0; i < typer.NumField(); i++ {
			field := typer.Field(i)
			// only support anonymous field
			if field.Anonymous {
				continue
			}
			// only support struct tag `query:"key"`
			queryKey, ok := field.Tag.Lookup("query")
			if !ok {
				continue
			}
			// url query doesn't have this key
			if !queries.Has(queryKey) {
				continue
			}
			// try to set query value to struct field
			value := queries.Get(queryKey)
			switch field.Type.Kind() {
			case reflect.String:
				valueT.Field(i).SetString(value)
			case reflect.Bool:
				if value == "true" {
					valueT.Field(i).SetBool(true)
				} else {
					valueT.Field(i).SetBool(false)
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				// failed to parse will return 0
				intx, _ := strconv.ParseInt(value, 10, 64)
				valueT.Field(i).SetInt(intx)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				// failed to parse will return 0
				uintx, _ := strconv.ParseUint(value, 10, 64)
				valueT.Field(i).SetUint(uintx)
			case reflect.Float32, reflect.Float64:
				// failed to parse will return 0
				float, _ := strconv.ParseFloat(value, 64)
				valueT.Field(i).SetFloat(float)
			default:
				continue
			}
		}
		ctx.Map(valueT.Interface())
	}
}
