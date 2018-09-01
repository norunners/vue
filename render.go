package vue

import (
	"bytes"
	"fmt"
	"github.com/albrow/vdom"
	"github.com/cbroglie/mustache"
	"github.com/gowasm/go-js-dom"
	"github.com/tdewolff/minify"
	minhtml "github.com/tdewolff/minify/html"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"strings"
)

const (
	v     = "v-"
	vBind = "v-bind"
	vIf   = "v-if"
)

type renderer struct {
	tmpl []byte
	el   dom.Element
	tree *vdom.Tree
	flag *html.Node
}

// newRenderer creates a new renderer.
func newRenderer(el string, tmpl []byte) *renderer {
	element := dom.GetWindow().Document().QuerySelector(el)

	minifier := minify.New()
	minifier.Add("text/html", &minhtml.Minifier{KeepEndTags: true})
	tmpl, err := minifier.Bytes("text/html", tmpl)
	must(err)

	return &renderer{el: element, tmpl: tmpl, tree: &vdom.Tree{}, flag: &html.Node{}}
}

// render executes the template with the given data and applies it to the dom element.
func (renderer *renderer) render(data map[string]interface{}) {
	buf := bytes.NewBuffer(renderer.tmpl)
	nodes, err := html.ParseFragment(buf, &html.Node{
		Type:     html.ElementNode,
		Data:     "div",
		DataAtom: atom.Div,
	})
	must(err)

	node := renderer.renderNode(nodes[0], data)

	buf = bytes.NewBuffer(nil)
	err = html.Render(buf, node)
	must(err)

	tmpl, err := mustache.ParseString(buf.String())
	must(err)

	buf.Reset()
	err = tmpl.FRender(buf, data)
	must(err)

	tree, err := vdom.Parse(buf.Bytes())
	must(err)

	patches, err := vdom.Diff(renderer.tree, tree)
	must(err)

	patches.Patch(renderer.el)
	renderer.tree = tree
}

// renderNode recursive traverses the html tree and renders the nodes.
func (renderer *renderer) renderNode(node *html.Node, data map[string]interface{}) *html.Node {
	// Render attributes.
	for i, attr := range node.Attr {
		if strings.HasPrefix(attr.Key, v) {
			deleteAttr(node, i)
			node = renderer.renderAttr(node, attr, data)
			// The flag signals that the tree structure was modified.
			// The next sibling of flag is the node to render next.
			if node == renderer.flag {
				return node
			}
		}
	}

	// Render children.
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		child = renderer.renderNode(child, data)
	}
	// The flag must be removed if used, this preserves the expected html structure.
	// The flag node intentionally fails to render.
	if node == renderer.flag.Parent {
		node.RemoveChild(renderer.flag)
	}

	return node
}

// renderAttr renders the given vue attribute.
func (renderer *renderer) renderAttr(node *html.Node, attr html.Attribute, data map[string]interface{}) *html.Node {
	vals := strings.Split(attr.Key, ":")
	dir, part := vals[0], ""
	if len(vals) > 1 {
		part = vals[1]
	}
	switch dir {
	case vIf:
		node = renderer.renderAttrIf(node, attr.Val, data)
	case vBind:
		renderAttrBind(node, part, attr.Val)
	default:
		must(fmt.Errorf("unknown vue attribute: %v", dir))
	}
	return node
}

// renderAttrIf renders the vue if attribute.
func (renderer *renderer) renderAttrIf(node *html.Node, field string, data map[string]interface{}) *html.Node {
	if value, ok := data[field]; ok {
		if val, ok := value.(bool); ok && val {
			return node
		}
	}
	node.Parent.InsertBefore(renderer.flag, node)
	node.Parent.RemoveChild(node)
	return renderer.flag
}

// renderAttrBind renders the vue bind attribute.
func renderAttrBind(node *html.Node, key, value string) {
	node.Attr = append(node.Attr, html.Attribute{Key: key, Val: fmt.Sprintf("{{ %v }}", value)})
}

// deleteAttr deletes the attribute of the node at the index.
// Attribute order is not preserved.
func deleteAttr(node *html.Node, i int) {
	n := len(node.Attr) - 1
	node.Attr[i] = node.Attr[n]
	node.Attr = node.Attr[:n]
}
