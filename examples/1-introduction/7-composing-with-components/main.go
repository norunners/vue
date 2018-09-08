package main

import (
	"github.com/norunners/vue"
)

const tmpl = `
<div>
  <ol>
    <todo-item
      v-for="Item in Todos"
      v-bind:Todo="Item">
    </todo-item>
  </ol>
</div>
`

type Data struct {
	Todos []Todo
}

type Todo struct {
	Text string
}

func main() {
	data := &Data{
		Todos: []Todo{
			{Text: "Vegetables"},
			{Text: "Cheese"},
			{Text: "Whatever else humans are supposed to eat"},
		},
	}

	vue.New(
		vue.El("#app"),
		vue.Template(tmpl),
		vue.Data(data),
		vue.Sub("todo-item", vue.Component(
			vue.Props("Todo"),
			vue.Template("<li>{{ Todo.Text }}</li>"),
		)),
	)

	select {}
}
