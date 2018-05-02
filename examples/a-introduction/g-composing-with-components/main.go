// Package main is an example of composing with components.
// See guide: https://vuejs.org/v2/guide/index.html#Composing-with-Components for details.
package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/norunners/vue"
)

type data struct {
	*js.Object

	groceryList []todo `js:"groceryList"`
}

type todo struct {
	*js.Object

	text string `js:"text"`
}

func main() {
	data := &data{Object: newObject()}
	data.groceryList = []todo{
		{Object: newObject()},
		{Object: newObject()},
		{Object: newObject()},
	}
	data.groceryList[0].text = "Vegetables"
	data.groceryList[1].text = "Cheese"
	data.groceryList[2].text = "Whatever else humans are supposed to eat"

	app := vue.New(
		vue.El("#app"),
		vue.DataValue(data),
		vue.Components(
			vue.Field("todo-item",
				vue.Make(
					vue.PropsList("todo"),
					vue.Template("<li>{{ todo.text }}</li>"),
				),
			),
		),
	)
	js.Global.Set("app", app)
}

func newObject(args ...interface{}) *js.Object {
	return js.Global.Get("Object").New(args...)
}
