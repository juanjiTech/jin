package jin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNameOfFunction(t *testing.T) {
	f := func() {}
	assert.Equal(t, "github.com/juanjiTech/jin.TestNameOfFunction.func1", nameOfFunction(f))
	assert.Equal(t, "nil", nameOfFunction(nil))
	assert.Equal(t, "nil", nameOfFunction(HandlerFunc(nil)))
}

func TestLastChar(t *testing.T) {
	assert.Equal(t, uint8('a'), lastChar("a"))
	assert.Equal(t, uint8('b'), lastChar("ab"))
	assert.Equal(t, uint8(0), lastChar(""))
}

func TestJoinPaths(t *testing.T) {
	assert.Equal(t, "/a/b", joinPaths("/a", "/b"))
	assert.Equal(t, "/a/b", joinPaths("/a/", "/b"))
	assert.Equal(t, "/a/b", joinPaths("/a", "b"))
	assert.Equal(t, "/a/b", joinPaths("/a/", "b"))
	assert.Equal(t, "a/b", joinPaths("a", "/b"))
	assert.Equal(t, "a/b", joinPaths("a/", "/b"))
	assert.Equal(t, "a/b", joinPaths("a", "b"))
	assert.Equal(t, "a/b", joinPaths("a/", "b"))
	assert.Equal(t, "/a", joinPaths("/a", ""))
	assert.Equal(t, "/b", joinPaths("", "/b"))
	assert.Equal(t, "/", joinPaths("/", ""))
	assert.Equal(t, "/", joinPaths("", "/"))
	assert.Equal(t, "", joinPaths("", ""))
	assert.Equal(t, "/a/b/", joinPaths("/a", "b/"))
}

func TestResolveAddress(t *testing.T) {
	// Test with environment variable
	t.Setenv("PORT", "8081")
	addr := resolveAddress([]string{})
	assert.Equal(t, ":8081", addr)

	// Test with passed address
	addr = resolveAddress([]string{":8082"})
	assert.Equal(t, ":8082", addr)

	// Test default
	t.Setenv("PORT", "")
	addr = resolveAddress([]string{})
	assert.Equal(t, ":8080", addr)
}

func TestOrdinalize(t *testing.T) {
	assert.Equal(t, "1st", ordinalize(1))
	assert.Equal(t, "2nd", ordinalize(2))
	assert.Equal(t, "3rd", ordinalize(3))
	assert.Equal(t, "4th", ordinalize(4))
	assert.Equal(t, "11th", ordinalize(11))
	assert.Equal(t, "12th", ordinalize(12))
	assert.Equal(t, "13th", ordinalize(13))
	assert.Equal(t, "21st", ordinalize(21))
	assert.Equal(t, "22nd", ordinalize(22))
	assert.Equal(t, "23rd", ordinalize(23))
	assert.Equal(t, "101st", ordinalize(101))
}
