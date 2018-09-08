package vue

import (
	"fmt"
	"github.com/fatih/structs"
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
// Props are not included in data.
func (vm *ViewModel) Data() interface{} {
	return vm.comp.data
}

// Get returns the data field value.
// Props are included to get.
func (vm *ViewModel) Get(field string) interface{} {
	data := vm.dataMap()
	value, ok := data[field]
	if !ok {
		must(fmt.Errorf("unknown data field: %s", field))
	}
	return value
}

// Set assigns the data field to the given value.
// Props are not included to set.
func (vm *ViewModel) Set(field string, value interface{}) {
	data := reflect.Indirect(reflect.ValueOf(vm.comp.data))
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

// dataMap returns the data plus props as a map.
func (vm *ViewModel) dataMap() map[string]interface{} {
	data := structs.Map(vm.comp.data)
	for field, prop := range vm.comp.props {
		data[field] = prop
	}
	return data
}
