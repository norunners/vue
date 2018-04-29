package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/norunners/vue"
)

type data struct {
	*js.Object

	message string `js:"message"`
}

func (data *data) reverseMessage() {
	r := []rune(data.message)
	n := len(r)
	for i := 0; i < n/2; i++ {
		r[i], r[n-i-1] = r[n-i-1], r[i]
	}
	data.message = string(r)
}

// main is an example of handling user input.
// See guide: https://vuejs.org/v2/guide/index.html#Handling-User-Input for details.
func main() {
	data := &data{Object: newObject()}
	data.message = "Hello World!"

	app := vue.New(
		vue.El("#app"),
		vue.DataValue(data),
		vue.Methods(
			vue.Field("reverseMessage", data.reverseMessage),
		),
	)

	js.Global.Set("app", app)
}

func newObject(args ...interface{}) *js.Object {
	return js.Global.Get("Object").New(args...)
}
