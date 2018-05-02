// Package main is an example of declarative rendering.
// See guide: https://vuejs.org/v2/guide/index.html#Declarative-Rendering for details.
package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/norunners/vue"
	"time"
)

type data struct {
	*js.Object

	message string `js:"message"`
}

func main() {
	data := &data{Object: newObject()}
	data.message = "You loaded this page on " + time.Now().String()

	app := vue.New(
		vue.El("#app"),
		vue.DataValue(data),
	)
	js.Global.Set("app", app)
}

func newObject(args ...interface{}) *js.Object {
	return js.Global.Get("Object").New(args...)
}
