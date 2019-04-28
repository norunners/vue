package main

import (
	"encoding/json"
	"github.com/norunners/vue"
	"net/http"
	"strings"
)

const tmpl = `
<div>
  <p>
    Ask a yes or no question:
    <input v-model="Question">
  </p>
  <p>{{ Answer }}</p>
</div>
`

type Data struct {
	Question string
	Answer   string
}

type yesno struct {
	Answer string `json:"answer"`
}

func Answer(vctx vue.Context, newQuestion, _ string) {
	if !strings.HasSuffix(newQuestion, "?") {
		vctx.Set("Answer", "Questions usually contain a question mark.")
		return
	}

	vctx.Go("AsyncAnswer")
}

func AsyncAnswer(vctx vue.Context) {
	data := vctx.Data().(*Data)
	res, err := http.Get("https://yesno.wtf/api")
	if err != nil {
		data.Answer = err.Error()
		return
	}
	defer res.Body.Close()

	dec := json.NewDecoder(res.Body)
	yesno := &yesno{}
	err = dec.Decode(yesno)
	if err != nil {
		data.Answer = err.Error()
		return
	}
	data.Answer = yesno.Answer
}

func main() {
	vue.New(
		vue.El("#app"),
		vue.Template(tmpl),
		vue.Data(&Data{Answer: "I cannot give you an answer until you ask a question!"}),
		vue.Watch("Question", Answer),
		vue.Methods(AsyncAnswer),
	)

	select {}
}
