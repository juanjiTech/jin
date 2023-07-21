package jin

import (
	"flag"
	"io"
	"os"
)

// EnvJinMode indicates environment name for Jin mode.
const EnvJinMode = "JIN_MODE"

const (
	// DebugMode indicates Jin mode is debug.
	DebugMode = "debug"
	// ReleaseMode indicates Jin mode is release.
	ReleaseMode = "release"
	// TestMode indicates Jin mode is test.
	TestMode = "test"
)

const (
	debugCode = iota
	releaseCode
	testCode
)

// DefaultWriter is the default io.Writer used by Jin for debug output and
// middleware output like Logger() or Recovery().
// Note that both Logger and Recovery provides custom ways to configure their
// output io.Writer.
// To support coloring in Windows use:
//
//	import "github.com/mattn/go-colorable"
//	jin.DefaultWriter = colorable.NewColorableStdout()
var DefaultWriter io.Writer = os.Stdout

// DefaultErrorWriter is the default io.Writer used by Jin to debug errors
var DefaultErrorWriter io.Writer = os.Stderr

var (
	jinMode  = debugCode
	modeName = DebugMode
)

func init() {
	mode := os.Getenv(EnvJinMode)
	SetMode(mode)
}

// SetMode sets gin mode according to input string.
func SetMode(value string) {
	if value == "" {
		if flag.Lookup("test.v") != nil {
			value = TestMode
		} else {
			value = DebugMode
		}
	}

	switch value {
	case DebugMode:
		jinMode = debugCode
	case ReleaseMode:
		jinMode = releaseCode
	case TestMode:
		jinMode = testCode
	default:
		panic("jin mode unknown: " + value + " (available mode: debug release test)")
	}

	modeName = value
}

// Mode returns current Jin mode.
func Mode() string {
	return modeName
}
