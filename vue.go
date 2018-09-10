// Package vue is the progressive framework for wasm applications.
package vue

// ViewModel is a vue view model, e.g. VM.
type ViewModel struct {
	comp      *Comp
	tmpl      *template
	renderer  *renderer
	rendered  bool
	data      map[string]interface{}
	callbacks map[string]struct{}
}

// New creates a new view model from the given options.
func New(options ...Option) *ViewModel {
	comp := Component(options...)
	return newViewModel(comp)
}

// newViewModel creates a new view model from the given component.
func newViewModel(comp *Comp) *ViewModel {
	tmpl := newTemplate(comp)
	renderer := newRenderer(comp)
	callbacks := make(map[string]struct{}, 0)

	vm := &ViewModel{comp: comp, tmpl: tmpl, renderer: renderer, callbacks: callbacks}
	// The root view model satisfies callback which is passed down to subcomponents.
	if comp.callback == nil {
		comp.callback = vm
	}
	vm.render()
	return vm
}

// must panics on errors.
func must(err error) {
	if err != nil {
		panic(err)
	}
}
