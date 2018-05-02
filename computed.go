package vue

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/norunners/vert"
)

// Computed is the vue computed option for components.
// See API: https://vuejs.org/v2/api/#computed for details.
func Computed(options ...Option) Option {
	computed := Make(options...)
	return Field("computed", computed)
}

// Get is the vue get option for computed.
func Get(function interface{}) Option {
	return Field("get", function)
}

// Set is the vue set option for computed.
func Set(function interface{}) Option {
	fn := vert.New(function)
	set := func(object *js.Object) {
		_, err := fn.Call(object)
		must(err)
	}
	return Field("set", set)
}
