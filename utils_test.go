package jin

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNameOfFunction(t *testing.T) {
	f := func() {}
	assert.Equal(t, "github.com/juanjiTech/jin.TestNameOfFunction.func1", nameOfFunction(f))
	assert.Equal(t, "nil", nameOfFunction(nil))
	assert.Equal(t, "nil", nameOfFunction(HandlerFunc(nil)))
}

func TestGetFuncReturnTypes(t *testing.T) {
	// Test function with no return values
	t.Run("NoReturn", func(t *testing.T) {
		fn := func() {}
		info := GetFuncReturnTypes(fn)
		assert.Equal(t, 0, info.NumOut)
		assert.Empty(t, info.OutTypes)
	})

	// Test function with single return value
	t.Run("SingleReturn", func(t *testing.T) {
		fn := func() int { return 0 }
		info := GetFuncReturnTypes(fn)
		assert.Equal(t, 1, info.NumOut)
		assert.Equal(t, reflect.TypeOf(0), info.OutTypes[0])
	})

	// Test function with multiple return values
	t.Run("MultipleReturns", func(t *testing.T) {
		fn := func() (int, string, error) { return 0, "", nil }
		info := GetFuncReturnTypes(fn)
		assert.Equal(t, 3, info.NumOut)
		assert.Equal(t, reflect.TypeOf(0), info.OutTypes[0])
		assert.Equal(t, reflect.TypeOf(""), info.OutTypes[1])
		assert.Equal(t, reflect.TypeOf((*error)(nil)).Elem(), info.OutTypes[2])
	})

	// Test function with struct return
	t.Run("StructReturn", func(t *testing.T) {
		type User struct {
			Name string
			Age  int
		}
		fn := func() User { return User{} }
		info := GetFuncReturnTypes(fn)
		assert.Equal(t, 1, info.NumOut)
		assert.Equal(t, reflect.TypeOf(User{}), info.OutTypes[0])
	})

	// Test function with pointer return
	t.Run("PointerReturn", func(t *testing.T) {
		fn := func() *int { return nil }
		info := GetFuncReturnTypes(fn)
		assert.Equal(t, 1, info.NumOut)
		assert.Equal(t, reflect.TypeOf((*int)(nil)), info.OutTypes[0])
	})

	// Test panic on non-function
	t.Run("PanicOnNonFunction", func(t *testing.T) {
		assert.Panics(t, func() {
			GetFuncReturnTypes(42)
		})
		assert.Panics(t, func() {
			GetFuncReturnTypes("not a function")
		})
		assert.Panics(t, func() {
			GetFuncReturnTypes(nil)
		})
	})
}

func TestCallFuncAndGetReturns(t *testing.T) {
	// Test function with no return values
	t.Run("NoReturn", func(t *testing.T) {
		called := false
		fn := func() { called = true }
		values, types := CallFuncAndGetReturns(fn)
		assert.True(t, called)
		assert.Empty(t, values)
		assert.Empty(t, types)
	})

	// Test function with single return value
	t.Run("SingleReturn", func(t *testing.T) {
		fn := func() int { return 42 }
		values, types := CallFuncAndGetReturns(fn)
		assert.Equal(t, 1, len(values))
		assert.Equal(t, 42, values[0])
		assert.Equal(t, reflect.TypeOf(0), types[0])
	})

	// Test function with multiple return values
	t.Run("MultipleReturns", func(t *testing.T) {
		fn := func() (int, string, bool) { return 100, "hello", true }
		values, types := CallFuncAndGetReturns(fn)
		assert.Equal(t, 3, len(values))
		assert.Equal(t, 100, values[0])
		assert.Equal(t, "hello", values[1])
		assert.Equal(t, true, values[2])
		assert.Equal(t, reflect.TypeOf(0), types[0])
		assert.Equal(t, reflect.TypeOf(""), types[1])
		assert.Equal(t, reflect.TypeOf(true), types[2])
	})

	// Test function with arguments
	t.Run("WithArguments", func(t *testing.T) {
		fn := func(a int, b string) (int, string) { return a * 2, b + "!" }
		values, types := CallFuncAndGetReturns(fn, 5, "test")
		assert.Equal(t, 2, len(values))
		assert.Equal(t, 10, values[0])
		assert.Equal(t, "test!", values[1])
		assert.Equal(t, reflect.TypeOf(0), types[0])
		assert.Equal(t, reflect.TypeOf(""), types[1])
	})

	// Test function with error return
	t.Run("ErrorReturn", func(t *testing.T) {
		fn := func(shouldErr bool) (string, error) {
			if shouldErr {
				return "", errors.New("test error")
			}
			return "success", nil
		}
		// Test without error
		values, types := CallFuncAndGetReturns(fn, false)
		assert.Equal(t, "success", values[0])
		assert.Nil(t, values[1])

		// Test with error
		values, types = CallFuncAndGetReturns(fn, true)
		assert.Equal(t, "", values[0])
		assert.NotNil(t, values[1])
		assert.Equal(t, "test error", values[1].(error).Error())
		assert.Equal(t, reflect.TypeOf((*error)(nil)).Elem(), types[1])
	})

	// Test panic on non-function
	t.Run("PanicOnNonFunction", func(t *testing.T) {
		assert.Panics(t, func() {
			CallFuncAndGetReturns(42)
		})
	})
}
