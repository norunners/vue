package vue

import (
	"github.com/gowasm/go-js-dom"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
)

// Option uses the option pattern for components.
type Option func(*Comp)

// El is the element option for components.
// The root element of a component is query selected from the value, e.g. #app or body.
func El(el string) Option {
	return func(comp *Comp) {
		comp.el = dom.GetWindow().Document().QuerySelector(el)
	}
}

// Template is the template option for components.
// The template uses the mustache syntax for rendering.
// The template must have a single root element.
func Template(tmpl string) Option {
	return func(comp *Comp) {
		minifier := minify.New()
		minifier.Add("text/html", &html.Minifier{KeepEndTags: true})
		tmpl, err := minifier.Bytes("text/html", []byte(tmpl))
		must(err)
		comp.tmpl = tmpl
	}
}

// Data is the data option for components.
// The scope of the data is within the component.
// Data must be a pointer to be mutable by methods.
func Data(data interface{}) Option {
	return func(comp *Comp) {
		comp.data = data
	}
}

// Sub is the subcomponent option for components.
func Sub(element string, sub *Comp) Option {
	return func(comp *Comp) {
		comp.subs[element] = sub
	}
}

// Props is the props option for subcomponents.
func Props(props ...string) Option {
	return func(sub *Comp) {
		for _, prop := range props {
			sub.props[prop] = nil
		}
	}
}
