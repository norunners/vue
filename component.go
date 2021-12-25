// Package vue is the progressive framework for wasm applications.
package vue

import (
	"github.com/gowasm/go-js-dom"
)

// Component is a vue component.
type Component struct {
	el      string
	tmpl    string
	data    interface{}
	methods map[string]func(*Component)
	root    renderer
}

// New creates a new component from the given options.
func New(options ...Option) *Component {
	methods := make(map[string]func(*Component), 0)
	comp := &Component{methods: methods}
	for _, option := range options {
		option(comp)
	}

	el := dom.GetWindow().Document().QuerySelector(comp.el)
	el.SetInnerHTML(comp.tmpl)

	comp.root = comp.newRenderer(el)
	comp.root.render(comp.data)

	return comp
}

// Data returns the data for the component.
func (comp *Component) Data() interface{} {
	return comp.data
}

// Call calls the method for the given name and renders the data.
func (comp *Component) Call(name string) {
	if function, ok := comp.methods[name]; ok {
		function(comp)
		comp.root.render(comp.data)
	}
}

// must panics on errors.
func must(err error) {
	if err != nil {
		panic(err)
	}
}
