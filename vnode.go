package vue

import (
	"fmt"
	"github.com/gowasm/go-js-dom"
	"golang.org/x/net/html"
	"strings"
	"syscall/js"
)

var document dom.Document

type vnode struct {
	parent, firstChild, lastChild, prevSibling, nextSibling *vnode

	attrs map[string]string
	typ   html.NodeType
	data  string

	node dom.Node
}

func init() {
	doc := js.Global().Get("document")
	if doc == js.Undefined() || doc == js.Null() {
		panic("failed to initialize document")
	}
	document = dom.WrapDocument(doc)
}

// newNode recursively creates a new virtual node from the dom node.
func newNode(node dom.Node) *vnode {
	if node == nil {
		return nil
	}
	vnode := &vnode{}
	switch n := node.(type) {
	case dom.Element:
		vnode.typ = html.ElementNode
		vnode.data = strings.ToLower(n.TagName())

		attrs := n.Attributes()
		vnode.attrs = make(map[string]string, len(attrs))
		for key, val := range attrs {
			vnode.attrs[key] = val
		}

		for child := node.FirstChild(); child != nil; child = child.NextSibling() {
			vnode.append(newNode(child))
		}
	case dom.Text:
		vnode.typ = html.TextNode
		vnode.data = n.TextContent()
	default:
		must(fmt.Errorf("unknown dom node type: %d", n.NodeType()))
	}
	// Set the dom node last prevents dom calls during creation.
	vnode.node = node
	return vnode
}

// render recursively renders the virtual node.
func (dst *vnode) render(src *html.Node) {
	for dstChild, srcChild := dst.firstChild, src.FirstChild; dstChild != nil || srcChild != nil; {
		switch {
		case dstChild == nil:
			dst.append(createNode(srcChild))
		case srcChild == nil:
			dst.remove(dstChild)
		case dstChild.typ != srcChild.Type:
			dst.replace(createNode(srcChild), dstChild)
		default:
			switch srcChild.Type {
			case html.ElementNode:
				if dstChild.data != srcChild.Data {
					dst.replace(createNode(srcChild), dstChild)
				} else {
					dstChild.renderAttributes(dstChild.attrs)
					dstChild.render(srcChild)
				}
			case html.TextNode:
				if dstChild.data != srcChild.Data {
					dstChild.setText(srcChild.Data)
				}
			default:
				must(fmt.Errorf("unknown html node type: %v", srcChild.Type))
			}
		}
		if dstChild != nil {
			dstChild = dstChild.nextSibling
		}
		if srcChild != nil {
			srcChild = srcChild.NextSibling
		}
	}
}

// createNode recursively creates a virtual node from the html node.
// createNode recursively creates a virtual node from the html node.
func createNode(node *html.Node) *vnode {
	vnode := &vnode{typ: node.Type, data: node.Data}
	switch node.Type {
	case html.ElementNode:
		vnode.node = document.CreateElement(node.Data)
		vnode.attrs = make(map[string]string, len(node.Attr))
		for _, attr := range node.Attr {
			vnode.setAttr(attr.Key, attr.Val)
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			vnode.append(createNode(child))
		}
	case html.TextNode:
		vnode.node = document.CreateTextNode(node.Data)
	default:
		must(fmt.Errorf("unknown node type: %v", node.Type))
	}
	return vnode
}

// renderAttributes renders the attributes.
func (vnode *vnode) renderAttributes(attrs map[string]string) {
	keys := make(map[string]struct{}, len(vnode.attrs)+len(attrs))
	for key := range vnode.attrs {
		keys[key] = struct{}{}
	}
	for key := range attrs {
		keys[key] = struct{}{}
	}

	for key := range keys {
		if srcVal, ok := attrs[key]; ok {
			if dstVal, ok := vnode.attrs[key]; !ok || dstVal != srcVal {
				vnode.setAttr(key, srcVal)
			}
		} else {
			vnode.remAttr(key)
		}
	}
}

// setAttr sets an attribute of the element.
func (vnode *vnode) setAttr(key, val string) {
	vnode.attrs[key] = val
	if vnode.node != nil {
		vnode.node.(dom.Element).SetAttribute(key, val)
	}
}

// remAttr removes an attribute from the element.
func (vnode *vnode) remAttr(key string) {
	delete(vnode.attrs, key)
	if vnode.node != nil {
		vnode.node.(dom.Element).RemoveAttribute(key)
	}
}

// setText sets the content of the text.
func (vnode *vnode) setText(content string) {
	vnode.data = content
	if vnode.node != nil {
		vnode.node.SetTextContent(content)
	}
}

// append appends the child to the node.
func (vnode *vnode) append(child *vnode) {
	prev := vnode.lastChild
	if prev == nil {
		vnode.firstChild = child
	} else {
		prev.nextSibling = child
	}
	vnode.lastChild = child
	child.parent = vnode
	child.prevSibling = prev

	if vnode.node != nil {
		vnode.node.AppendChild(child.node)
	}
}

// replace replaces a child with a new child.
func (vnode *vnode) replace(newChild, oldChild *vnode) {
	prev, next := oldChild.prevSibling, oldChild.nextSibling
	if prev == nil {
		vnode.firstChild = newChild
	} else {
		prev.nextSibling = newChild
	}
	if next == nil {
		vnode.lastChild = newChild
	} else {
		next.prevSibling = newChild
	}
	newChild.parent = vnode
	newChild.prevSibling = prev
	newChild.nextSibling = next

	if vnode.node != nil {
		vnode.node.ReplaceChild(newChild.node, oldChild.node)
	}
}

// remove removes a child from the node.
func (vnode *vnode) remove(child *vnode) {
	if vnode.firstChild == child {
		vnode.firstChild = child.nextSibling
	}
	if child.nextSibling != nil {
		child.nextSibling.prevSibling = child.prevSibling
	}
	if vnode.lastChild == child {
		vnode.lastChild = child.prevSibling
	}
	if child.prevSibling != nil {
		child.prevSibling.nextSibling = child.nextSibling
	}

	if vnode.node != nil {
		vnode.node.RemoveChild(child.node)
	}
}
