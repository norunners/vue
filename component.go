// Package vue is the progressive framework for wasm applications.
package vue

import "github.com/fatih/structs"

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

// render calls the renderer with the prepared data.
func (comp *Component) render() {
	data := structs.Map(comp.data)
	comp.renderer.render(data)
}

// Call calls the method for the given name and calls render.
func (comp *Component) Call(name string) {
	if function, ok := comp.methods[name]; ok {
		function(comp)
		comp.render()
	}
}

// must panics on errors.
func must(err error) {
	if err != nil {
		panic(err)
	}
}
