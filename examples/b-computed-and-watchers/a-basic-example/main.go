// Package main is an example of a basic computed property.
// See guide: https://vuejs.org/v2/guide/computed.html#Basic-Example for details.
package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/norunners/vue"
)

type data struct {
	*js.Object

	message string `js:"message"`
}

func (data *data) reversedMessage() string {
	r := []rune(data.message)
	n := len(r)
	for i := 0; i < n/2; i++ {
		r[i], r[n-i-1] = r[n-i-1], r[i]
	}
	return string(r)
}

func main() {
	data := &data{Object: newObject()}
	data.message = "Hello World!"

	app := vue.New(
		vue.El("#app"),
		vue.DataValue(data),
		vue.Computed(
			vue.Field("reversedMessage", data.reversedMessage),
		),
	)
	js.Global.Set("app", app)
}

func newObject(args ...interface{}) *js.Object {
	return js.Global.Get("Object").New(args...)
}
