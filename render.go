package vue

import (
	"fmt"
	"github.com/fatih/structs"
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
	vm.state = structs.Map(vm.data)
	vm.mapProps()
	vm.mapComputed()
}

// mapProps maps props to state.
func (vm *ViewModel) mapProps() {
	for field, prop := range vm.props {
		vm.state[field] = prop
	}
}

// mapComputed maps computed to state.
func (vm *ViewModel) mapComputed() {
	for computed, function := range vm.comp.computed {
		if _, ok := vm.state[computed]; !ok {
			vm.state[computed] = function(vm)
		}
	}
}
