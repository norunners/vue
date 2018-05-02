// Package main is an example of watchers.
// See guide: https://vuejs.org/v2/guide/computed.html#Watchers for details.
package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/norunners/vue"
	"math/rand"
	"strings"
	"time"
)

type data struct {
	*js.Object

	question string `js:"question"`
	answer   string `js:"answer"`
}

func (data *data) questions(question, oldQuestion string) {
	data.answer = "Waiting for you to stop typing..."
	if !strings.HasSuffix(question, "?") {
		data.answer = "Questions usually contain a question mark. ;-)"
		return
	}

	time.AfterFunc(500*time.Millisecond, func() {
		if rand.Intn(2) == 0 {
			data.answer = "No"
		} else {
			data.answer = "Yes"
		}
	})
}

func main() {
	data := &data{Object: newObject()}
	data.question = ""
	data.answer = "I cannot give you an answer until you ask a question!"

	app := vue.New(
		vue.El("#app"),
		vue.DataValue(data),
		vue.Watch(
			vue.Watcher("question", data.questions),
		),
	)
	js.Global.Set("app", app)
}

func newObject(args ...interface{}) *js.Object {
	return js.Global.Get("Object").New(args...)
}
