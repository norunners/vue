// Package vue is the progressive framework for wasm applications.
package vue

import (
	"github.com/fatih/structs"
	"reflect"
)

// Component is a vue component.
type Component struct {
	el       string
	tmpl     []byte
	data     interface{}
	methods  map[string]func(Context)
	renderer *renderer
}

// New creates a new component from the given options.
func New(options ...Option) *Component {
	methods := make(map[string]func(Context), 0)
	comp := &Component{methods: methods}
	for _, option := range options {
		option(comp)
	}

	cbs := newCallbacks(comp)
	comp.renderer = newRenderer(comp.el, comp.tmpl, cbs)
	comp.render()
	comp.tmpl = nil

	return comp
}

// Data returns the data for the component.
func (comp *Component) Data() interface{} {
	return comp.data
}

// Get returns the data field value.
func (comp *Component) Get(field string) interface{} {
	data := reflect.Indirect(reflect.ValueOf(comp.data))
	return data.FieldByName(field).Interface()
}

// Set assigns the data field to the given value.
func (comp *Component) Set(field string, value interface{}) {
	data := reflect.Indirect(reflect.ValueOf(comp.data))
	val := reflect.Indirect(data.FieldByName(field))
	val.Set(reflect.Indirect(reflect.ValueOf(value)))
}

// Call calls the given method then calls render.
func (comp *Component) Call(method string) {
	if function, ok := comp.methods[method]; ok {
		function(comp)
		comp.render()
	}
}

// render calls the renderer with the prepared data.
func (comp *Component) render() {
	data := structs.Map(comp.data)
	comp.renderer.render(data)
}

// must panics on errors.
func must(err error) {
	if err != nil {
		panic(err)
	}
}
