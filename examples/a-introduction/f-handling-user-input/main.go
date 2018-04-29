package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/norunners/vue"
)

type data struct {
	*js.Object

	message string `js:"message"`
}

// main is an example of handling user input.
// See guide: https://vuejs.org/v2/guide/index.html#Handling-User-Input for details.
func main() {
	data := &data{Object: newObject()}
	data.message = "Hello World!"

	app := vue.New(
		vue.El("#app"),
		vue.DataValue(data),
	)

	js.Global.Set("app", app)
}

func newObject(args ...interface{}) *js.Object {
	return js.Global.Get("Object").New(args...)
}
