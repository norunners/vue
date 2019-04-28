package vue

import (
	"fmt"
	"reflect"
)

// Context is received by functions to interact with the component.
type Context interface {
	Data() interface{}
	Get(field string) interface{}
	Set(field string, value interface{})
	Go(method string, args ...interface{})
	Emit(event string, args ...interface{})
}

// Data returns the data for the component.
// Props and computed are excluded from data.
func (vm *ViewModel) Data() interface{} {
	return vm.data.Interface()
}

// Get returns the data field value.
// Props and computed are included to get.
// Computed may be calculated as needed.
func (vm *ViewModel) Get(field string) interface{} {
	if value, ok := vm.state[field]; ok {
		return value
	}

	function, ok := vm.comp.computed[field]
	if !ok {
		must(fmt.Errorf("unknown data field: %s", field))
	}
	value := vm.compute(function)
	vm.mapField(field, value)
	return value
}

// Set assigns the data field to the given value.
// Props and computed are excluded to set.
func (vm *ViewModel) Set(field string, value interface{}) {
	data := reflect.Indirect(vm.data)
	oldVal := reflect.Indirect(data.FieldByName(field))
	newVal := reflect.Indirect(reflect.ValueOf(value))

	oldVal.Set(newVal)
	vm.mapField(field, value)
}

// Go asynchronously calls the given method with optional arguments.
// Blocking functions must be called asynchronously.
func (vm *ViewModel) Go(method string, args ...interface{}) {
	values := make([]reflect.Value, 0, len(args))
	for _, arg := range args {
		values = append(values, reflect.ValueOf(arg))
	}
	go vm.call(method, values)
}

// Emit dispatches the given event with optional arguments.
func (vm *ViewModel) Emit(event string, args ...interface{}) {
	vm.bus.pub(event, "", args)
}

// call calls the given method with optional values then calls render.
func (vm *ViewModel) call(method string, values []reflect.Value) {
	if function, ok := vm.comp.methods[method]; ok {
		values = append([]reflect.Value{reflect.ValueOf(vm)}, values...)
		function.Call(values)
		vm.render()
	}
}

// compute calls the given function and returns the first element.
func (vm *ViewModel) compute(function reflect.Value) interface{} {
	values := []reflect.Value{reflect.ValueOf(vm)}
	rets := function.Call(values)
	return rets[0].Interface()
}
