package vue

import (
	"fmt"
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
	_, ok := vm.callbacks[typ]
	if ok {
		return
	}
	vm.comp.el.AddEventListener(typ, false, cb)
	vm.callbacks[typ] = struct{}{}
}

// vModel is the vue model event callback.
func (vm *ViewModel) vModel(event dom.Event) {
	typ := event.Type()
	field, ok := event.Target().Attributes()[typ]
	if !ok {
		must(fmt.Errorf("unknown event type: %s", typ))
	}

	value := event.Target().Underlying().Get("value").String()
	vm.Set(field, value)
	vm.render()
}

// vOn is the vue on event callback.
func (vm *ViewModel) vOn(event dom.Event) {
	typ := event.Type()
	method, ok := findMethod(event.Target(), typ)
	if !ok {
		must(fmt.Errorf("unknown event type: %s", typ))
	}

	vm.Call(method)
}

// findMethod finds the method name from the event type by searching up the dom tree.
func findMethod(elem dom.Element, typ string) (string, bool) {
	if elem == nil {
		return "", false
	}
	if method, ok := elem.Attributes()[typ]; ok {
		return method, true
	}
	return findMethod(elem.ParentElement(), typ)
}
