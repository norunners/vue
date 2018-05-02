// Package main is an example of conditionals and loops.
// See guide: https://vuejs.org/v2/guide/index.html#Conditionals-and-Loops for details.
package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/norunners/vue"
)

type data struct {
	*js.Object

	todos []todo `js:"todos"`
}

type todo struct {
	*js.Object

	text string `js:"text"`
}

func main() {
	data := &data{Object: newObject()}
	data.todos = []todo{
		{Object: newObject()},
		{Object: newObject()},
		{Object: newObject()},
	}
	data.todos[0].text = "Learn GopherJS"
	data.todos[1].text = "Learn Vue"
	data.todos[2].text = "Build something awesome"

	app := vue.New(
		vue.El("#app"),
		vue.DataValue(data),
	)
	js.Global.Set("app", app)
}

func newObject(args ...interface{}) *js.Object {
	return js.Global.Get("Object").New(args...)
}
