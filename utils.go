package jin

import (
	"math"
	"os"
	"path"
	"reflect"
	"runtime"
	"strconv"
)

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
