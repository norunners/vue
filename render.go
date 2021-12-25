package vue

import (
	"fmt"
	"github.com/cbroglie/mustache"
	"github.com/gowasm/go-js-dom"
	"strings"
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
			tmpl, err := mustache.ParseString(content)
			must(err)
			return &text{node: node, tmpl: tmpl}
		}
	default:
		must(fmt.Errorf("unknown node: %v with type: %T", node, node))
	}
	return nil
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

// text represents a text element.
type text struct {
	node dom.Text
	tmpl *mustache.Template
	prev string
}

// render satisfies renderer for text elements.
func (text *text) render(data ...interface{}) {
	next, err := text.tmpl.Render(data...)
	must(err)
	if text.prev != next {
		text.node.SetTextContent(next)
		text.prev = next
	}
}
