package vue

import (
	"reflect"
	"runtime"
	"strings"
)

// Option uses the option pattern for components.
type Option func(*Comp)

// El is the element option for components.
// The root element of a component is query selected from the value, e.g. #app or body.
func El(el string) Option {
	return func(comp *Comp) {
		comp.el = document.QuerySelector(el)
	}
}

// Template is the template option for components.
// The template uses the mustache syntax for rendering.
// The template must have a single root element.
func Template(tmpl string) Option {
	return func(comp *Comp) {
		comp.tmpl = tmpl
	}
}

// Data is the data option for components.
// The scope of the data is within the component.
// Data must be a pointer to be mutable by methods.
func Data(data interface{}) Option {
	return func(comp *Comp) {
		comp.data = data
	}
}

// Methods is the methods option for components.
// The given functions are registered as methods for the component.
func Methods(functions ...func(context Context)) Option {
	return func(comp *Comp) {
		for _, function := range functions {
			name := funcName(function)
			comp.methods[name] = function
		}
	}
}

// Computed is the computed option for components.
// The given functions are registered as computed properties for the component.
func Computed(functions ...func(Context) interface{}) Option {
	return func(comp *Comp) {
		for _, function := range functions {
			name := funcName(function)
			comp.computed[name] = function
		}
	}
}

// Sub is the subcomponent option for components.
func Sub(element string, sub *Comp) Option {
	return func(comp *Comp) {
		comp.subs[element] = sub
	}
}

// Props is the props option for subcomponents.
func Props(props ...string) Option {
	return func(sub *Comp) {
		for _, prop := range props {
			sub.props[prop] = nil
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
