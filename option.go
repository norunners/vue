package vue

// Option uses the option pattern for components.
type Option func(*Component)

// El is the element option for components.
// The root element of a component is query selected from the value, e.g. #app or body.
func El(el string) Option {
	return func(comp *Component) {
		comp.el = el
	}
}

// Template is the template option for components.
// The template uses the mustache syntax for rendering.
func Template(tmpl string) Option {
	return func(comp *Component) {
		comp.tmpl = []byte(tmpl)
	}
}

// Data is the data option for components.
// The scope of the data is within the component.
// Data must be a pointer to be mutated by methods.
func Data(data interface{}) Option {
	return func(comp *Component) {
		comp.data = data
	}
}
