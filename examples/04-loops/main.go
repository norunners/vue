package main

import (
	"github.com/norunners/vue"
	"time"
)

const tmpl = `
<ol>
  <li v-for="Todo in Todos">
    {{ Todo.Text }}
  </li>
</ol>
`

type Data struct {
	Todos []Todo
}

type Todo struct {
	Text string
}

func Add(vctx vue.Context) {
	data := vctx.Data().(*Data)
	data.Todos = append(data.Todos, Todo{"Build something wasm!"})
}

func main() {
	data := &Data{
		Todos: []Todo{
			{Text: "Learn Go"},
			{Text: "Learn Vue"},
		},
	}

	vm := vue.New(
		vue.El("#app"),
		vue.Template(tmpl),
		vue.Data(data),
		vue.Methods(Add),
	)

	time.AfterFunc(time.Second, func() {
		vm.Go("Add")
	})
	select {}
}
