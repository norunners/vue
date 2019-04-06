// Package vue is the progressive framework for wasm applications.
package vue

import (
	"fmt"
	"reflect"
)

// Comp is a vue component.
type Comp struct {
	el       string
	tmpl     string
	data     interface{}
	methods  map[string]func(Context)
	computed map[string]func(Context) interface{}
	props    map[string]struct{}
	subs     map[string]*Comp
	isSub    bool
}

// Component creates a new component from the given options.
func Component(options ...Option) *Comp {
	methods := make(map[string]func(Context), 0)
	computed := make(map[string]func(Context) interface{}, 0)
	props := make(map[string]struct{}, 0)
	subs := make(map[string]*Comp, 0)

	comp := &Comp{
		data:     struct{}{},
		methods:  methods,
		computed: computed,
		props:    props,
		subs:     subs,
	}
	for _, option := range options {
		option(comp)
	}
	return comp
}

// newData creates new data from the function.
// Without a function the data of the component is returned.
func (comp *Comp) newData() interface{} {
	value := reflect.ValueOf(comp.data)
	if value.Type().Kind() != reflect.Func {
		return value.Interface()
	}
	rets := value.Call(nil)
	if n := len(rets); n != 1 {
		must(fmt.Errorf("invalid return length: %d", n))
	}
	return rets[0].Interface()
}
