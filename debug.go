package jin

import (
	"fmt"
	"strings"
)

func IsDebugging() bool {
	return jinMode == debugCode
}

// DebugPrintRouteFunc indicates debug log output format.
var DebugPrintRouteFunc func(httpMethod, absolutePath, handlerName string, nuHandlers int)

func debugPrintRoute(httpMethod, absolutePath string, handlers HandlersChain) {
	if IsDebugging() {
		nuHandlers := len(handlers)
		handlerName := nameOfFunction(handlers.Last())
		if DebugPrintRouteFunc == nil {
			debugPrint("%-6s %-25s --> %s (%d handlers)\n", httpMethod, absolutePath, handlerName, nuHandlers)
		} else {
			DebugPrintRouteFunc(httpMethod, absolutePath, handlerName, nuHandlers)
		}
	}
}

func debugPrint(format string, values ...any) {
	if IsDebugging() {
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}
		_, _ = fmt.Fprintf(DefaultWriter, "[JIN-debug] "+format, values...)
	}
}

func debugPrintWARNINGNew() {
	debugPrint(`[WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export JIN_MODE=release
 - using code:	jin.SetMode(jin.ReleaseMode)

`)
}

func debugPrintError(err error) {
	if err != nil && IsDebugging() {
		_, _ = fmt.Fprintf(DefaultErrorWriter, "[JIN-debug] [ERROR] %v\n", err)
	}
}
