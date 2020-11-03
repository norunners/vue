package vue

import (
	"fmt"
	"syscall/js"

	"golang.org/x/net/html"
	dom "honnef.co/go/js/dom/v2"
)

var document = dom.WrapDocument(js.Global().Get("document"))

type vnode struct {
	parent, firstChild, lastChild, prevSibling, nextSibling *vnode

	attrs map[string]string
	typ   html.NodeType
	data  string

	node dom.Node
}

// newNode creates a virtual node by query selecting the given element.
func newNode(el string) *vnode {
	node := document.QuerySelector(el)
	return &vnode{attrs: node.Attributes(), node: node}
}

// newSubNode creates a virtual subcomponent node from the given template.
func newSubNode(tmpl string) *vnode {
	node := parseNode(tmpl)
	var ok bool
	if node, ok = firstElement(node); !ok {
		must(fmt.Errorf("failed to find first element from template: %s", tmpl))
	}
	return createElement(node)
}

// createElement creates a virtual node element without children nor attributes.
func createElement(node *html.Node) *vnode {
	el := document.CreateElement(node.Data)
	attrs := make(map[string]string, len(node.Attr))
	return &vnode{
		typ:   node.Type,
		data:  node.Data,
		attrs: attrs,
		node:  el,
	}
}

// createNode recursively creates a virtual node from the html node.
func createNode(node *html.Node, subs subs) *vnode {
	vnode := &vnode{typ: node.Type, data: node.Data}
	switch node.Type {
	case html.ElementNode:
		if subNode, ok := subs.vnode(node.Data); ok {
			subNode.renderAttributes(node.Attr)
			return subNode
		} else {
			vnode.node = document.CreateElement(node.Data)
			vnode.attrs = make(map[string]string, len(node.Attr))
			for _, attr := range node.Attr {
				vnode.setAttr(attr.Key, attr.Val)
			}
			for child := node.FirstChild; child != nil; child = child.NextSibling {
				vnode.append(createNode(child, subs))
			}
		}
	case html.TextNode:
		vnode.node = document.CreateTextNode(node.Data)
	default:
		must(fmt.Errorf("unknown node type: %v", node.Type))
	}
	return vnode
}

// render recursively renders the virtual node.
func (dst *vnode) render(src *html.Node, subs subs) {
	for dstChild, srcChild := dst.firstChild, src.FirstChild; dstChild != nil || srcChild != nil; {
		switch {
		case dstChild == nil:
			dst.append(createNode(srcChild, subs))
		case srcChild == nil:
			dst.remove(dstChild)
		case dstChild.typ != srcChild.Type:
			dst.replace(createNode(srcChild, subs), dstChild)
		default:
			switch srcChild.Type {
			case html.ElementNode:
				if subNode, ok := subs.vnode(srcChild.Data); ok {
					subNode.renderAttributes(srcChild.Attr)
					dst.replace(subNode, dstChild)
				} else if dstChild.data != srcChild.Data {
					dst.replace(createNode(srcChild, subs), dstChild)
				} else {
					dstChild.renderAttributes(srcChild.Attr)
					dstChild.render(srcChild, subs)
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

// renderAttributes renders the attributes.
func (vnode *vnode) renderAttributes(attrs []html.Attribute) {
	keys := make(map[string]struct{}, len(vnode.attrs)+len(attrs))
	srcAttrs := make(map[string]string, len(attrs))
	for _, attr := range attrs {
		keys[attr.Key] = struct{}{}
		srcAttrs[attr.Key] = attr.Val
	}
	for key := range vnode.attrs {
		keys[key] = struct{}{}
	}

	for key := range keys {
		if srcVal, ok := srcAttrs[key]; ok {
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
		if key == "value" {
			vnode.node.Underlying().Set(key, val)
		}
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
