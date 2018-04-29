package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/norunners/vue"
)

type data struct {
	*js.Object

	seen bool `js:"seen"`
}

// main is an example of conditionals and loops.
// See guide: https://vuejs.org/v2/guide/index.html#Conditionals-and-Loops for details.
func main() {
	data := &data{Object: newObject()}
	data.seen = true

	app := vue.New(
		vue.El("#app"),
		vue.DataValue(data),
	)

	js.Global.Set("app", app)
}

func newObject(args ...interface{}) *js.Object {
	return js.Global.Get("Object").New(args...)
}
