package jin

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSetMode(t *testing.T) {
	_ = os.Setenv(EnvJinMode, TestMode)
	mode := os.Getenv(EnvJinMode)
	SetMode(mode)

	assert.Equal(t, testCode, jinMode)
	assert.Equal(t, TestMode, Mode())
	_ = os.Unsetenv(EnvJinMode)

	SetMode("")
	assert.Equal(t, testCode, jinMode)
	assert.Equal(t, TestMode, Mode())

	tmp := flag.CommandLine
	flag.CommandLine = flag.NewFlagSet("", flag.ContinueOnError)
	SetMode("")
	assert.Equal(t, debugCode, jinMode)
	assert.Equal(t, DebugMode, Mode())
	flag.CommandLine = tmp

	SetMode(DebugMode)
	assert.Equal(t, debugCode, jinMode)
	assert.Equal(t, DebugMode, Mode())

	SetMode(ReleaseMode)
	assert.Equal(t, releaseCode, jinMode)
	assert.Equal(t, ReleaseMode, Mode())

	SetMode(TestMode)
	assert.Equal(t, testCode, jinMode)
	assert.Equal(t, TestMode, Mode())

	assert.Panics(t, func() { SetMode("unknown") })
}
