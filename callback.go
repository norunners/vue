package vue

import (
	"github.com/gowasm/go-js-dom"
)

// callback interacts with event listeners on the root element.
// The callback is passed down to subcomponents.
type callback interface {
	addEventListener(typ string, cb func(dom.Event))
	vModel(event dom.Event)
	vOn(event dom.Event)
	render()
}

// addEventListener adds the callback to the element as an event listener unless the type was previously added.
func (vm *ViewModel) addEventListener(typ string, cb func(dom.Event)) {
	if _, ok := vm.callbacks[typ]; ok {
		return
	}
	vm.comp.el.AddEventListener(typ, cb, false)
	vm.callbacks[typ] = struct{}{}
}

// vModel is the vue model event callback.
func (vm *ViewModel) vModel(event dom.Event) {
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
	method, ok := findAttrValue(event.Target(), event.Type())
	if !ok {
		return
	}
	vm.Call(method)
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
