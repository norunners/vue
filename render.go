package vue

import (
	"fmt"
	"github.com/cbroglie/mustache"
	"github.com/gowasm/go-js-dom"
	"strings"
)

const (
	v     = "v-"
	vBind = "v-bind"
)

// renderer renders the given data.
type renderer interface {
	render(data ...interface{})
}

// newRenderer creates a new renderer from the given dom node
// by traversing the dom tree recursively.
func (comp *Component) newRenderer(src dom.Node) renderer {
	switch node := src.(type) {
	case dom.Element:
		children := make([]renderer, 0)
		for key, val := range node.Attributes() {
			if strings.HasPrefix(key, v) {
				attr := newAttr(node, key, val)
				children = append(children, attr)
			}
		}
		for _, child := range node.ChildNodes() {
			if renderer := comp.newRenderer(child); renderer != nil {
				children = append(children, renderer)
			}
		}
		if len(children) > 0 {
			return &parent{children: children}
		}
	case dom.Text:
		if content := strings.TrimSpace(node.TextContent()); content != "" {
			return newText(node, content)
		}
	default:
		must(fmt.Errorf("unknown node: %v with type: %T", node, node))
	}
	return nil
}

// newText creates a new renderer for a text node.
func newText(text dom.Text, content string) renderer {
	tmpl, err := mustache.ParseString(content)
	must(err)

	update := func(next string) {
		text.SetTextContent(next)
	}
	return &updater{tmpl: tmpl, update: update}
}

// newText creates a new renderer for a vue attribute.
func newAttr(el dom.Element, key, value string) renderer {
	vals := strings.Split(key, ":")
	key, typ := vals[0], ""
	if len(vals) > 1 {
		typ = vals[1]
	}
	switch key {
	case vBind:
		return newBindAttr(el, typ, value)
	default:
		must(fmt.Errorf("unknown vue directive: %v", key))
	}
	return nil
}

// newBindAttr creates a new renderer for a vue bind attribute.
func newBindAttr(el dom.Element, key, value string) renderer {
	tmpl, err := mustache.ParseString(fmt.Sprintf("{{ %v }}", value))
	must(err)

	update := func(next string) {
		el.SetAttribute(key, next)
	}
	return &updater{tmpl: tmpl, update: update}
}

// parent represents an element with children.
type parent struct {
	children []renderer
}

// render satisfies renderer for parent elements.
func (parent *parent) render(data ...interface{}) {
	for _, child := range parent.children {
		child.render(data...)
	}
}

// updater renders templates and calls update on changes.
type updater struct {
	tmpl   *mustache.Template
	update func(next string)
	prev   string
}

// render satisfies renderer for updater.
func (up *updater) render(data ...interface{}) {
	next, err := up.tmpl.Render(data...)
	must(err)
	if up.prev != next {
		up.update(next)
		up.prev = next
	}
}
