// Package vue is the progressive framework for wasm applications.
package vue

import "syscall/js"

// ViewModel is a vue view model, e.g. VM.
type ViewModel struct {
	comp  *Comp
	vnode *vnode
	data  interface{}
	state map[string]interface{}
	funcs map[string]js.Func
	props map[string]interface{}
	subs  subs

	index int
}

// New creates a new view model from the given options.
func New(options ...Option) *ViewModel {
	comp := Component(options...)
	return newViewModel(comp, nil)
}

// newViewModel creates a new view model from the given component with props.
func newViewModel(comp *Comp, props map[string]interface{}) *ViewModel {
	var vnode *vnode
	if comp.isSub {
		vnode = newSubNode(comp.tmpl)
	} else {
		vnode = newNode(comp.el)
	}
	data := comp.newData()
	funcs := make(map[string]js.Func, 0)
	subs := newSubs(comp.subs)

	vm := &ViewModel{
		comp:  comp,
		vnode: vnode,
		data:  data,
		funcs: funcs,
		props: props,
		subs:  subs,
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
