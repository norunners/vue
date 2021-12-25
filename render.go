package vue

import (
	"bytes"
	"github.com/norunners/vdom"
	"golang.org/x/net/html"
)

// render calls the renderer with the prepared data.
// Subcomponents calls the callback to render from the root element.
func (vm *ViewModel) render() {
	if vm.comp.isSub {
		if vm.rendered {
			vm.comp.callback.render()
		}
		return
	}

	vm.mapData()
	b := vm.tmpl.execute(vm.data)
	vm.renderer.render(b)
}

// subRender renders the subcomponent into a node.
func (vm *ViewModel) subRender() *html.Node {
	vm.mapData()
	b := vm.tmpl.execute(vm.data)
	reader := bytes.NewReader(b)
	nodes := parse(reader)
	vm.rendered = true
	return nodes[0]
}

// renderer interacts with the virtual dom.
type renderer struct {
	comp *Comp
	tree *vdom.Tree
}

// newRenderer creates a new renderer.
func newRenderer(comp *Comp) *renderer {
	return &renderer{comp: comp, tree: &vdom.Tree{}}
}

// render uses the virtual dom to render the given bytes for the root dom element.
func (renderer *renderer) render(b []byte) {
	tree, err := vdom.Parse(b)
	must(err)

	patches, err := vdom.Diff(renderer.tree, tree)
	must(err)

	err = patches.Patch(renderer.comp.el)
	must(err)

	renderer.tree = tree
}
