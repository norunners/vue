package main

import (
	"github.com/norunners/vue"
	"math/rand"
	"time"
)

const tmpl = `
<p v-bind:class="Class">
	Hello WebAssembly!
</p>
`

type Data struct {
	Class Class
}

type Class struct {
	Active     bool `css:"active"`
	TextDanger bool `css:"text-danger"`
}

func Change(vctx vue.Context) {
	data := vctx.Data().(*Data)
	data.Class.Active = rand.Intn(2) == 1
	data.Class.TextDanger = rand.Intn(2) == 1
}

func main() {
	vm := vue.New(
		vue.El("#app"),
		vue.Template(tmpl),
		vue.Data(&Data{}),
		vue.Methods(Change),
	)

	for tick := time.Tick(time.Second); ; {
		select {
		case <-tick:
			vm.Go("Change")
		}
	}
}
