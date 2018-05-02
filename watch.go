package vue

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/norunners/vert"
)

// Watch is the vue watch option for components.
// See API: https://vuejs.org/v2/api/#watch for details.
func Watch(options ...Option) Option {
	watch := Make(options...)
	return Field("watch", watch)
}

// Watcher is the watcher option for watch.
// See guide: https://vuejs.org/v2/guide/computed.html#Watchers for details.
func Watcher(key string, function interface{}) Option {
	fn := vert.New(function)
	watcher := func(object, oldObject *js.Object) {
		_, err := fn.Call(object, oldObject)
		must(err)
	}
	return Field(key, watcher)
}

// Handler is the vue handler option for watch.
func Handler(function interface{}) Option {
	return Watcher("handler", function)
}

// Deep is the vue deep option for watch.
func Deep() Option {
	return Field("deep", true)
}

// Immediate is the vue immediate option for watch.
func Immediate() Option {
	return Field("immediate", true)
}
