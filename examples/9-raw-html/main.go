package main

import (
	"github.com/norunners/vue"
)

const (
	tmpl = `
<div>
  <p>Using mustaches: {{{ RawHtml }}}</p>
  <p>Using v-html directive: <span v-html="RawHtml"></span></p>
</div>
`
	rawHtml = `
<span style="color: red">This should be red.</span>
`
)

type Data struct {
	RawHtml string
}

func main() {
	vue.New(
		vue.El("#app"),
		vue.Template(tmpl),
		vue.Data(Data{RawHtml: rawHtml}),
	)

	select {}
}
