package main

import (
	"github.com/norunners/vue"
	"time"
)

const tmpl = `
<span v-bind:title="Message">
  Hover your mouse over me for a few seconds
  to see my dynamically bound title!
</span>
`

type Data struct {
	Message string
}

func main() {
	vue.New(
		vue.El("#app"),
		vue.Template(tmpl),
		vue.Data(Data{Message: "You loaded this page on " + time.Now().Format(time.ANSIC)}),
	)

	select {}
}
