package vue

import (
	"reflect"
	"runtime"
	"strings"
)

// Methods is the methods option for components.
// The given functions are registered as methods for the component.
func Methods(functions ...func(context Context)) Option {
	return func(comp *Component) {
		for _, function := range functions {
			name := funcName(function)
			comp.methods[name] = function
		}
	}
}

// funcName returns the name of the given function.
func funcName(function interface{}) string {
	fn := reflect.ValueOf(function)
	name := runtime.FuncForPC(fn.Pointer()).Name()
	return stripMetadata(name)
}

// stripMetadata returns the function name without metadata.
func stripMetadata(name string) string {
	parts := strings.Split(name, ".")
	name = parts[len(parts)-1]
	return strings.TrimRight(name, "-fm")
}
