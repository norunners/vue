package vue

import (
	"reflect"
)

// bus contains subscriptions of events to methods.
type bus struct {
	parent *bus
	caller caller
	subs   map[string]map[string]struct{}
}

// caller calls a method with optional arguments.
type caller interface {
	call(method string, args []reflect.Value)
}

// newBus creates a new event bus.
func newBus(parent *bus, caller caller) *bus {
	subs := make(map[string]map[string]struct{}, 0)
	return &bus{parent: parent, caller: caller, subs: subs}
}

// pub publishes the event with the method and optional arguments.
// All event subscriber methods are called if the method is empty.
// When the component is not subscribed to the event, it propagates to parent components.
func (bus *bus) pub(event, method string, args []interface{}) {
	if bus == nil {
		return
	}

	methods, ok := bus.subs[event]
	if !ok {
		bus.parent.pub(event, method, args)
		return
	}

	values := make([]reflect.Value, 0, len(args))
	for _, arg := range args {
		values = append(values, reflect.ValueOf(arg))
	}

	if method != "" {
		bus.caller.call(method, values)
		return
	}

	for method := range methods {
		bus.caller.call(method, values)
	}
}

// sub subscribes the component to the event with the method.
func (bus *bus) sub(event, method string) {
	if methods, ok := bus.subs[event]; ok {
		methods[method] = struct{}{}
	} else {
		bus.subs[event] = map[string]struct{}{method: {}}
	}
}
