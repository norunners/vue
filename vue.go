// Package vue is the progressive framework for gopherjs applications.
package vue

import (
	"github.com/gopherjs/gopherjs/js"
)

// Map is the map representation of a js object.
type Map map[string]interface{}

// New instantiates a new Vue object from the given options.
// The Vue object is also known as a view model, e.g. vm.
// See guide: https://vuejs.org/v2/guide/instance.html for details.
func New(options ...Option) *js.Object {
	m := Make(options...)
	return js.Global.Get("Vue").New(m)
}

// Make makes a new map from the given options.
func Make(options ...Option) Map {
	m := make(Map, len(options))
	for _, option := range options {
		option(m)
	}
	return m
}

// must panics on errors.
func must(err error) {
	if err != nil {
		panic(err)
	}
}
