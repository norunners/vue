package vue

import (
	"fmt"
	"reflect"
)

// Context is received by methods to interact with the component.
type Context interface {
	Data() interface{}
	Get(field string) interface{}
	Set(field string, value interface{})
	Call(method string)
}

// Data returns the data for the component.
// Props and computed are excluded from data.
func (vm *ViewModel) Data() interface{} {
	return vm.data
}

// Get returns the data field value.
// Props and computed are included to get.
// Computed may be calculated as needed.
func (vm *ViewModel) Get(field string) interface{} {
	value, ok := vm.state[field]
	if !ok {
		function, ok := vm.comp.computed[field]
		if !ok {
			must(fmt.Errorf("unknown data field: %s", field))
		}

		value = function(vm)
		vm.state[field] = value
	}
	return value
}

// Set assigns the data field to the given value.
// Props and computed are excluded to set.
func (vm *ViewModel) Set(field string, value interface{}) {
	data := reflect.Indirect(reflect.ValueOf(vm.data))
	val := reflect.Indirect(data.FieldByName(field))
	val.Set(reflect.Indirect(reflect.ValueOf(value)))
}

// Call calls the given method then calls render.
func (vm *ViewModel) Call(method string) {
	if function, ok := vm.comp.methods[method]; ok {
		function(vm)
		vm.render()
	}
}
