// Package vue is the progressive framework for wasm applications.
package vue

import (
	"github.com/gowasm/go-js-dom"
)

// Comp is a vue component.
type Comp struct {
	el      dom.Element
	tmpl    []byte
	data    interface{}
	methods map[string]func(Context)
	subs    map[string]*Comp

	props    map[string]interface{}
	isSub    bool
	callback callback
}

// Component creates a new component from the given options.
func Component(options ...Option) *Comp {
	methods := make(map[string]func(Context), 0)
	subs := make(map[string]*Comp, 0)
	props := make(map[string]interface{}, 0)

	comp := &Comp{data: struct{}{}, methods: methods, subs: subs, props: props}
	for _, option := range options {
		option(comp)
	}
	return comp
}

// hasProp determines if a component has a prop.
// Returns false for nil components.
func (comp *Comp) hasProp(prop string) bool {
	if comp == nil {
		return false
	}
	_, ok := comp.props[prop]
	return ok
}

// newSub attempts to creates a new subcomponent.
// Returns false for unknown elements.
func (comp *Comp) newSub(element string) (*Comp, bool) {
	sub, ok := comp.subs[element]
	if !ok {
		return nil, false
	}
	sub.isSub = true
	sub.callback = comp.callback
	return sub, true
}
