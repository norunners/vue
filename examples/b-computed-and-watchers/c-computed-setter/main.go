// Package main is an example of a computed vs watched property.
// See guide: https://vuejs.org/v2/guide/computed.html#Computed-Setter for details.
package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/norunners/vue"
	"strings"
)

type data struct {
	*js.Object

	firstName string `js:"firstName"`
	lastName  string `js:"lastName"`
}

func (data *data) fullName() string {
	return data.firstName + " " + data.lastName
}

func (data *data) setFullName(fullName string) {
	names := strings.Split(fullName, " ")
	data.firstName = names[0]
	data.lastName = names[1]
}

func main() {
	data := &data{Object: newObject()}
	data.firstName = "Foo"
	data.lastName = "Bar"

	app := vue.New(
		vue.El("#app"),
		vue.DataValue(data),
		vue.Computed(
			vue.Field("fullName",
				vue.Make(
					vue.Get(data.fullName),
					vue.Set(data.setFullName),
				),
			),
		),
	)
	js.Global.Set("app", app)
}

func newObject(args ...interface{}) *js.Object {
	return js.Global.Get("Object").New(args...)
}
