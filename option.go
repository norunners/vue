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
		comp.el = el
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
// This option accepts either a function or a struct.
// The data function is expected to return a new data value.
// For example: func() *Type { return &Type{...} }
// Without a function the data is shared across components.
// The scope of the data is within the component.
// Data must be a pointer to be mutable by methods.
func Data(data interface{}) Option {
	return func(comp *Comp) {
		comp.data = data
	}
}

// Method is the method option for components.
// The given name and function is registered as a method for the component.
// The function is required to accept context and allows optional arguments.
// For example: func(ctx vue.Context) or func(ctx vue.Context, a1 Arg1, ..., ak ArgK)
func Method(name string, function interface{}) Option {
	return func(comp *Comp) {
		comp.methods[name] = reflect.ValueOf(function)
	}
}

// Methods is the methods option for components.
// The given functions are registered as methods for the component.
// The functions are required to accept context and allows optional arguments.
// For example: func(ctx vue.Context) or func(ctx vue.Context, a1 Arg1, ..., ak ArgK)
func Methods(functions ...interface{}) Option {
	return func(comp *Comp) {
		for _, function := range functions {
			fn := reflect.ValueOf(function)
			name := funcName(fn)
			comp.methods[name] = fn
		}
	}
}

// Computed is the computed option for components.
// The given name and function is registered as a computed property for the component.
// The function is required to accept context and return a value.
// For example: func(ctx vue.Context) Type
func Computed(name string, function interface{}) Option {
	return func(comp *Comp) {
		fn := reflect.ValueOf(function)
		comp.computed[name] = fn
	}
}

// Computeds is the computeds option for components.
// The given functions are registered as computed properties for the component.
// The functions are required to accept context and return a value.
// For example: func(ctx vue.Context) Type
func Computeds(functions ...interface{}) Option {
	return func(comp *Comp) {
		for _, function := range functions {
			fn := reflect.ValueOf(function)
			name := funcName(fn)
			comp.computed[name] = fn
		}
	}
}

// Watch is the watch option for components.
// The given function is registered as a watcher for the data field.
// All data fields are watchable, e.g. data, props and computed.
// The function is required to accept context and both the new and old values.
// For example: func(ctx vue.Context, newVal, oldVal Type)
func Watch(field string, function interface{}) Option {
	return func(comp *Comp) {
		fn := reflect.ValueOf(function)
		comp.watchers[field] = fn
	}
}

// Sub is the subcomponent option for components.
func Sub(element string, sub *Comp) Option {
	return func(comp *Comp) {
		sub.isSub = true
		comp.subs[element] = sub
	}
}

// Props is the props option for subcomponents.
func Props(props ...string) Option {
	return func(sub *Comp) {
		for _, prop := range props {
			sub.props[prop] = struct{}{}
		}
	}
}

// funcName returns the name of the given function.
func funcName(function reflect.Value) string {
	name := runtime.FuncForPC(function.Pointer()).Name()
	return stripMetadata(name)
}

// stripMetadata returns the function name without metadata.
func stripMetadata(name string) string {
	parts := strings.Split(name, ".")
	name = parts[len(parts)-1]
	return strings.TrimSuffix(name, "-fm")
}
