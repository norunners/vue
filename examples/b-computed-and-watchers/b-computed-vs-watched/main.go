// Package main is an example of a computed vs watched property.
// Note, this examples demonstrates an overuse of watches.
// See guide: https://vuejs.org/v2/guide/computed.html#Computed-vs-Watched-Property for details.
package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/norunners/vue"
)

type data struct {
	*js.Object

	firstName string `js:"firstName"`
	lastName  string `js:"lastName"`
	fullName  string `js:"fullName"`
}

func (data *data) watchFirstName(firstName string) {
	data.fullName = firstName + " " + data.lastName
}

func (data *data) watchLastName(lastName string) {
	data.fullName = data.firstName + " " + lastName
}

func main() {
	data := &data{Object: newObject()}
	data.firstName = "Foo"
	data.lastName = "Bar"
	data.fullName = "Foo Bar"

	app := vue.New(
		vue.El("#app"),
		vue.DataValue(data),
		vue.Watch(
			vue.Watcher("firstName", data.watchFirstName),
			vue.Watcher("lastName", data.watchLastName),
		),
	)
	js.Global.Set("app", app)
}

func newObject(args ...interface{}) *js.Object {
	return js.Global.Get("Object").New(args...)
}
