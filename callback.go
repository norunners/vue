package vue

import (
	"fmt"
	"github.com/gowasm/go-js-dom"
)

type callbacks struct {
	context Context
	set     map[string]struct{}
}

// newCallbacks is created with event callback methods.
func newCallbacks(context Context) *callbacks {
	set := make(map[string]struct{}, 0)
	return &callbacks{context: context, set: set}
}

// addCallback adds the callback to the element as an event listener unless the type was previously added.
func (cbs *callbacks) addCallback(el dom.Element, typ string, cb func(dom.Event)) {
	_, ok := cbs.set[typ]
	if ok {
		return
	}
	el.AddEventListener(typ, false, cb)
	cbs.set[typ] = struct{}{}
}

// vModel is the vue model event callback.
func (cbs *callbacks) vModel(event dom.Event) {
	typ := event.Type()
	field, ok := event.Target().Attributes()[typ]
	if !ok {
		must(fmt.Errorf("unknown event type: %s", typ))
	}

	value := event.Target().Underlying().Get("value").String()
	cbs.context.Set(field, value)
	cbs.context.render()
}

// vOn is the vue on event callback.
func (cbs *callbacks) vOn(event dom.Event) {
	typ := event.Type()
	method, ok := event.Target().Attributes()[typ]
	if !ok {
		must(fmt.Errorf("unknown event type: %s", typ))
	}

	cbs.context.Call(method)
}
