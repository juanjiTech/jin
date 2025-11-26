package jin

import (
	"fmt"
	"math"
	"os"
	"path"
	"reflect"
	"runtime"
	"strconv"
)

// FuncReturnInfo contains information about a function's return values
type FuncReturnInfo struct {
	NumOut  int            // Number of return values
	OutTypes []reflect.Type // Types of each return value
}

// GetFuncReturnTypes returns the return types of a function using reflection.
// It panics if the provided value is not a function.
func GetFuncReturnTypes(fn any) FuncReturnInfo {
	fnType := reflect.TypeOf(fn)
	if fnType == nil || fnType.Kind() != reflect.Func {
		panic(fmt.Sprintf("expected a function, got %T", fn))
	}

	numOut := fnType.NumOut()
	outTypes := make([]reflect.Type, numOut)
	for i := 0; i < numOut; i++ {
		outTypes[i] = fnType.Out(i)
	}

	return FuncReturnInfo{
		NumOut:   numOut,
		OutTypes: outTypes,
	}
}

// CallFuncAndGetReturns calls a function with the given arguments and returns
// the return values along with their types.
// It panics if fn is not a function or if the arguments don't match.
func CallFuncAndGetReturns(fn any, args ...any) (values []any, types []reflect.Type) {
	fnValue := reflect.ValueOf(fn)
	fnType := fnValue.Type()

	if fnType.Kind() != reflect.Func {
		panic(fmt.Sprintf("expected a function, got %T", fn))
	}

	// Prepare arguments
	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}

	// Call the function
	results := fnValue.Call(in)

	// Extract values and types
	values = make([]any, len(results))
	types = make([]reflect.Type, len(results))
	for i, result := range results {
		values[i] = result.Interface()
		types[i] = result.Type()
	}

	return values, types
}

func resolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			debugPrint("Environment variable PORT=\"%s\"", port)
			return ":" + port
		}
		debugPrint("Environment variable PORT is undefined. Using port :8080 by default")
		return ":8080"
	case 1:
		return addr[0]
	default:
		panic("too many parameters")
	}
}

func nameOfFunction(f any) string {
	if f == nil {
		return "nil"
	}
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}

func joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	if lastChar(relativePath) == '/' && lastChar(finalPath) != '/' {
		return finalPath + "/"
	}
	return finalPath
}

// ordinalize ordinalizes the number by adding the ordinal to the number.
func ordinalize(number int) string {
	abs := int(math.Abs(float64(number)))

	nstr := strconv.Itoa(number)
	i := abs % 100
	if i == 11 || i == 12 || i == 13 {
		return nstr + "th"
	}

	switch abs % 10 {
	case 1:
		return nstr + "st"
	case 2:
		return nstr + "nd"
	case 3:
		return nstr + "rd"
	default:
		return nstr + "th"
	}
}
