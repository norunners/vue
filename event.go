package vue

import (
	"strings"

	dom "honnef.co/go/js/dom/v2"
)

// addEventListener adds the callback to the element as an event listener unless the type was previously added.
func (vm *ViewModel) addEventListener(typ string, cb func(dom.Event)) {
	if _, ok := vm.funcs[typ]; ok {
		return
	}
	fn := vm.vnode.node.AddEventListener(typ, false, cb)
	vm.funcs[typ] = fn
}

// vModel is the vue model event callback.
func (vm *ViewModel) vModel(event dom.Event) {
	event.StopImmediatePropagation()

	target := event.Target()
	_, field, ok := findAttr(target, event.Type())
	if !ok {
		return
	}

	value := target.Underlying().Get("value").String()
	vm.Set(field, value)
	vm.render()
}

// vOn is the vue on event callback.
func (vm *ViewModel) vOn(event dom.Event) {
	event.StopImmediatePropagation()

	typ := event.Type()
	attrKey, method, ok := findAttr(event.Target(), typ)
	if !ok {
		return
	}

	modifiers := strings.TrimPrefix(attrKey, typ)
	modSet := modSet(modifiers)

	if keyEvent, ok := event.(*dom.KeyboardEvent); ok {
		key := keyEvent.Key()
		if _, ok := modSet[key]; !ok && len(modSet) > 0 {
			return
		}
	}

	vm.bus.pub(typ, method, nil)
}

// release removes all the event listeners.
func (vm *ViewModel) release() {
	for typ, fn := range vm.funcs {
		vm.vnode.node.RemoveEventListener(typ, false, fn)
		fn.Release()
	}
}

// findAttr finds the attribute from the given prefix by searching up the dom tree.
func findAttr(elem dom.Element, prefix string) (string, string, bool) {
	if elem == nil {
		return "", "", false
	}
	for attrKey, attrVal := range elem.Attributes() {
		if strings.HasPrefix(attrKey, prefix) {
			return attrKey, attrVal, true
		}
	}
	return findAttr(elem.ParentElement(), prefix)
}

// modSet converts modifiers to a set, includes title conversion.
// For example: hello.world -> {"Hello", "World"}
func modSet(modifiers string) map[string]struct{} {
	if modifiers == "" {
		return nil
	}
	mods := strings.Split(modifiers, ".")
	set := make(map[string]struct{}, len(mods))
	for _, mod := range mods {
		set[modTitle(mod)] = struct{}{}
	}
	return set
}

// modTitle converts modifiers to title style.
// For example: page-down -> PageDown
func modTitle(modifier string) string {
	mods := strings.Split(modifier, "-")
	sb := &strings.Builder{}
	for _, mod := range mods {
		sb.WriteString(strings.Title(mod))
	}
	return sb.String()
}
