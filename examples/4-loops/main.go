package main

import (
	"github.com/norunners/vue"
	"time"
)

const tmpl = `
<div>
  <ol>
    <li v-for="Todo in Todos">
      {{ Todo.Text }}
    </li>
  </ol>
</div>
`

type Data struct {
	Todos []Todo
}

type Todo struct {
	Text string
}

func Add(context vue.Context) {
	data := context.Data().(*Data)
	data.Todos = append(data.Todos, Todo{"Yikes!"})
}

func main() {
	data := &Data{
		Todos: []Todo{
			{Text: "Learn wasm"},
			{Text: "Learn vue"},
			{Text: "Build something awesome"},
		},
	}

	vm := vue.New(
		vue.El("#app"),
		vue.Template(tmpl),
		vue.Data(data),
		vue.Methods(Add),
	)

	time.AfterFunc(time.Second, func() {
		vm.Call("Add")
	})
	select {}
}
