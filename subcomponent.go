package vue

// subs maps elements to subcomponents
type subs map[string]*sub

// sub contains all the subcomponent instances for a component.
type sub struct {
	comp      *Comp
	index     int
	instances map[int]*instance
}

// instance contains a view model with props.
type instance struct {
	props map[string]interface{}
	vm    *ViewModel
}

// newSubs creates a new map of subcomponents.
func newSubs(comps map[string]*Comp) subs {
	subs := make(subs, len(comps))
	for element, comp := range comps {
		subs[element] = newSub(comp)
	}
	return subs
}

// newSub creates a new subcomponent.
func newSub(comp *Comp) *sub {
	instances := make(map[int]*instance, 0)
	return &sub{comp: comp, instances: instances}
}

// putProp puts the props in the subcomponent.
// Returns false if the element is not a subcomponent
// or the subcomponent is not expecting the prop.
func (subs subs) putProp(element, field string, data interface{}) bool {
	sub, ok := subs[element]
	if !ok {
		return false
	}
	return sub.putProp(field, data)
}

// putProp puts the props in the instance.
// Returns false if the subcomponent is not expecting the prop.
func (sub *sub) putProp(field string, data interface{}) bool {
	if _, ok := sub.comp.props[field]; !ok {
		return false
	}

	if inst, ok := sub.instances[sub.index]; ok {
		if inst.props == nil {
			inst.props = map[string]interface{}{field: data}
		} else {
			inst.props[field] = data
		}
	} else {
		sub.instances[sub.index] = &instance{props: map[string]interface{}{field: data}}
	}
	return true
}

// newInstance creates a new instance of the subcomponent with props.
// // Returns false if the element is not a subcomponent.
func (subs subs) newInstance(element string) bool {
	sub, ok := subs[element]
	if !ok {
		return false
	}
	return sub.newInstance()
}

// newInstance creates a new instance of the subcomponent with props.
func (sub *sub) newInstance() bool {
	if inst, ok := sub.instances[sub.index]; ok {
		if inst.vm == nil {
			inst.vm = newViewModel(sub.comp, inst.props)
		} else {
			inst.vm.props = inst.props
			inst.vm.render()
		}
	} else {
		vm := newViewModel(sub.comp, nil)
		sub.instances[sub.index] = &instance{vm: vm}
	}
	sub.index++
	return true
}

// vnode retrieves a virtual node of the subcomponent.
// // Returns false if the element is not a subcomponent.
func (subs subs) vnode(element string) (*vnode, bool) {
	sub, ok := subs[element]
	if !ok {
		return nil, false
	}
	vnode, ok := sub.vnode()
	return vnode, ok
}

// vnode retrieves a virtual node of the subcomponent.
func (sub *sub) vnode() (*vnode, bool) {
	inst, ok := sub.instances[sub.index]
	if !ok {
		return nil, false
	}
	sub.index++
	return inst.vm.vnode, true
}

// reset resets all subcomponents.
func (subs subs) reset() {
	for _, sub := range subs {
		sub.reset()
	}
}

// reset cleans up and unmounts unused subcomponent instances.
func (sub *sub) reset() {
	for i := sub.index; i < len(sub.instances); {
		sub.instances[i].vm.release()
		delete(sub.instances, i)
	}
	sub.index = 0
}
