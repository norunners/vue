// Package vue is the progressive framework for wasm applications.
package vue

import (
	"github.com/gowasm/go-js-dom"
)

// Component is a vue component.
type Component struct {
	el   string
	tmpl string
	data interface{}
	root renderer
}

// New creates a new component from the given options.
func New(options ...Option) *Component {
	comp := &Component{}
	for _, option := range options {
		option(comp)
	}

	el := dom.GetWindow().Document().QuerySelector(comp.el)
	el.SetInnerHTML(comp.tmpl)

	comp.root = comp.newRenderer(el)
	comp.root.render(comp.data)

	return comp
}

// must panics on errors.
func must(err error) {
	if err != nil {
		panic(err)
	}
}
