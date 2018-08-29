package main

import (
	"github.com/norunners/vue"
	"time"
)

const tmpl = `
<div>
  <span v-if="Seen">Now you see me</span>
</div>
`

type Data struct {
	Seen bool
}

func ToggleSeen(context vue.Context) {
	data := context.Data().(*Data)
	data.Seen = !data.Seen
}

func main() {
	comp := vue.New(
		vue.El("#app"),
		vue.Template(tmpl),
		vue.Data(&Data{Seen: true}),
		vue.Methods(ToggleSeen),
	)

	for ticker := time.NewTicker(time.Second); ; {
		select {
		case <-ticker.C:
			comp.Call("ToggleSeen")
		}
	}
}
