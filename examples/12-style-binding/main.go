package main

import (
	"fmt"
	"github.com/norunners/vue"
	"math/rand"
	"time"
)

const tmpl = `
<p v-bind:style="Style">
	Hello WebAssembly!
</p>
`

type Data struct {
	r, g, b int
	px      int
}

type Styles struct {
	Color    string `css:"color"`
	FontSize string `css:"font-size"`
}

func Style(vctx vue.Context) *Styles {
	data := vctx.Data().(*Data)
	hex := fmt.Sprintf("#%02x%02x%02x", data.r, data.g, data.b)
	size := fmt.Sprintf("%dpx", data.px)
	return &Styles{
		Color:    hex,
		FontSize: size,
	}
}

func Change(vctx vue.Context) {
	data := vctx.Data().(*Data)
	data.r = int(rand.Float32() * 0xff)
	data.g = int(rand.Float32() * 0xff)
	data.b = int(rand.Float32() * 0xff)
	data.px = 8 + (data.px-7)%64
}

func main() {
	vm := vue.New(
		vue.El("#app"),
		vue.Template(tmpl),
		vue.Data(&Data{px: 8}),
		vue.Computeds(Style),
		vue.Methods(Change),
	)

	for tick := time.Tick(200 * time.Millisecond); ; {
		select {
		case <-tick:
			vm.Go("Change")
		}
	}
}
