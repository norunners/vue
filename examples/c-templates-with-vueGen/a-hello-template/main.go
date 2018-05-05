// Package main is an example of the template option.
package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/norunners/vue"
)

const appTmpl = "<div>{{ message }}</div>"

type data struct {
	*js.Object

	message string `js:"message"`
}

func main() {
	data := &data{Object: newObject()}
	data.message = "Hello World!"

	app := vue.New(
		vue.El("#app"),
		vue.Template(appTmpl),
		vue.DataValue(data),
	)
	js.Global.Set("app", app)
}

func newObject(args ...interface{}) *js.Object {
	return js.Global.Get("Object").New(args...)
}
