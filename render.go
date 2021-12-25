package vue

import (
	"fmt"
	"reflect"
)

// render executes and renders the prepared state.
func (vm *ViewModel) render() {
	vm.mapState()
	node := vm.execute(vm.state)
	vm.subs.reset()
	if vm.comp.isSub {
		var ok bool
		if node, ok = firstElement(node); !ok {
			must(fmt.Errorf("failed to find first element from node: %s", node.Data))
		}
	}
	vm.vnode.render(node, vm.subs)
	vm.subs.reset()
}

// mapData creates a map of state from data, props and computed.
func (vm *ViewModel) mapState() {
	elem := reflect.Indirect(vm.data)
	typ := elem.Type()
	n := elem.NumField()
	vm.state = make(map[string]interface{}, n)
	for i := 0; i < n; i++ {
		field := elem.Field(i)
		if field.CanInterface() {
			name := typ.Field(i).Name
			value := field.Interface()
			vm.mapField(name, value)
		}
	}
	vm.mapProps()
	vm.mapComputed()
}

// mapProps maps props to state.
func (vm *ViewModel) mapProps() {
	for field, prop := range vm.props {
		vm.mapField(field, prop)
	}
}

// mapComputed maps computed to state.
func (vm *ViewModel) mapComputed() {
	for computed, function := range vm.comp.computed {
		if _, ok := vm.state[computed]; !ok {
			value := vm.compute(function)
			vm.mapField(computed, value)
		}
	}
}

// mapField maps a field to state.
// Watchers are called on field changes.
func (vm *ViewModel) mapField(field string, value interface{}) {
	oldField, ok := vm.state[field]
	vm.state[field] = value
	if !ok {
		return
	}

	if watcher, ok := vm.comp.watchers[field]; ok {
		newVal := reflect.ValueOf(value)
		oldVal := reflect.ValueOf(oldField)
		if reflect.DeepEqual(newVal, oldVal) {
			return
		}
		values := append([]reflect.Value{reflect.ValueOf(vm)}, newVal, oldVal)
		watcher.Call(values)
	}
}
