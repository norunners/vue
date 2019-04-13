package vue

import (
	"github.com/gowasm/go-js-dom"
)

// addEventListener adds the callback to the element as an event listener unless the type was previously added.
func (vm *ViewModel) addEventListener(typ string, cb func(dom.Event)) {
	if _, ok := vm.funcs[typ]; ok {
		return
	}
	fn := vm.vnode.node.AddEventListener(typ, cb, false)
	vm.funcs[typ] = fn
}

// vModel is the vue model event callback.
func (vm *ViewModel) vModel(event dom.Event) {
	event.StopImmediatePropagation()
	field, ok := findAttrValue(event.Target(), event.Type())
	if !ok {
		return
	}
	value := event.Target().Underlying().Get("value").String()
	vm.Set(field, value)
	vm.render()
}

// vOn is the vue on event callback.
func (vm *ViewModel) vOn(event dom.Event) {
	event.StopImmediatePropagation()
	method, ok := findAttrValue(event.Target(), event.Type())
	if !ok {
		return
	}
	vm.bus.pub(event.Type(), method, nil)
}

// release removes all the event listeners.
func (vm *ViewModel) release() {
	for typ, fn := range vm.funcs {
		vm.vnode.node.RemoveEventListener(typ, fn, false)
	}
}

// findAttrValue finds the attribute value from the given key by searching up the dom tree.
func findAttrValue(elem dom.Element, key string) (string, bool) {
	if elem == nil {
		return "", false
	}
	if method, ok := elem.Attributes()[key]; ok {
		return method, true
	}
	return findAttrValue(elem.ParentElement(), key)
}
