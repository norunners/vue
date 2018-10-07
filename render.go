package vue

import (
	"golang.org/x/net/html"
)

// render renders the prepared data.
// Subcomponents use the callback to render the root element.
func (vm *ViewModel) render() {
	if vm.comp.isSub {
		if vm.executed {
			vm.comp.callback.render()
		}
		return
	}

	vm.mapData()
	node := vm.tmpl.execute(vm.data)
	vm.vnode.render(node)
}

// executeSub executes the subcomponent into a node.
func (vm *ViewModel) executeSub() *html.Node {
	vm.mapData()
	node := vm.tmpl.execute(vm.data)
	vm.executed = true
	return node
}
