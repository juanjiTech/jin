package binding

import (
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/juanjiTech/jin"
)

var DefaultQueryFormat = "PREFIX.CURRENT"

func recursiveSetting(typer reflect.Type, queries url.Values, prefix string) (value reflect.Value) {
	if typer.Kind() == reflect.Ptr {
		typer = typer.Elem()
	}
	valueT := reflect.New(typer).Elem()

	for i := 0; i < typer.NumField(); i++ {
		field := typer.Field(i)
		// only support non-anonymous field
		if field.Anonymous {
			continue
		}
		// only support struct tag `query:"key"`
		queryKey, ok := field.Tag.Lookup("query")
		if !ok {
			continue
		}

		if prefix != "" {
			queryKey = strings.ReplaceAll(strings.ReplaceAll(DefaultQueryFormat, "PREFIX", prefix), "CURRENT", queryKey)
		}
		// url query doesn't have this key
		// if !queries.Has(queryKey) {
		//	continue
		// }
		// try to set query value to struct field
		switch field.Type.Kind() {
		case reflect.Struct:
			valueT.Field(i).Set(recursiveSetting(field.Type, queries, queryKey))
		case reflect.Slice:
			slicesType := field.Type.Elem()
			sliceValue := reflect.MakeSlice(field.Type, 0, 0)
			queryValue := queries.Get(queryKey)
			if queryValue == "" {
				continue
			}
			for _, v := range strings.Split(queryValue, ",") {
				newValue := reflect.New(slicesType).Elem()
				switch slicesType.Kind() {
				case reflect.String:
					newValue.SetString(v)
				case reflect.Bool:
					if v == "true" || v == "1" || v == "True" {
						newValue.SetBool(true)
					} else {
						newValue.SetBool(false)
					}
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					// failed to parse will return 0
					intx, _ := strconv.ParseInt(v, 10, 64)
					newValue.SetInt(intx)
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					// failed to parse will return 0
					uintx, _ := strconv.ParseUint(v, 10, 64)
					newValue.SetUint(uintx)
				case reflect.Float32, reflect.Float64:
					// failed to parse will return 0
					float, _ := strconv.ParseFloat(v, 64)
					newValue.SetFloat(float)
				}
				sliceValue = reflect.Append(sliceValue, newValue)
			}
			valueT.Field(i).Set(sliceValue)
		case reflect.String:
			value := queries.Get(queryKey)
			valueT.Field(i).SetString(value)
		case reflect.Bool:
			value := queries.Get(queryKey)
			if value == "true" || value == "1" || value == "True" {
				valueT.Field(i).SetBool(true)
			} else {
				valueT.Field(i).SetBool(false)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value := queries.Get(queryKey)
			// failed to parse will return 0
			intx, _ := strconv.ParseInt(value, 10, 64)
			valueT.Field(i).SetInt(intx)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			value := queries.Get(queryKey)
			// failed to parse will return 0
			uintx, _ := strconv.ParseUint(value, 10, 64)
			valueT.Field(i).SetUint(uintx)
		case reflect.Float32, reflect.Float64:
			value := queries.Get(queryKey)
			// failed to parse will return 0
			float, _ := strconv.ParseFloat(value, 64)
			valueT.Field(i).SetFloat(float)
		default:
			continue
		}
	}
	return valueT
}

func Query[T any](model T) jin.HandlerFunc {
	_ = model // to avoid unused variable warning
	typer := reflect.TypeOf(model)
	// only support struct
	if typer.Kind() == reflect.Ptr {
		typer = typer.Elem()
	}
	if typer.Kind() != reflect.Struct {
		panic("model must be a struct")
	}
	return func(ctx *jin.Context) {
		queries := ctx.Request.URL.Query()
		valueT := recursiveSetting(typer, queries, "")
		ctx.Map(valueT.Interface())
	}
}
