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
	"io"
	"reflect"
	"strings"
)

const (
	v      = "v-"
	vBind  = "v-bind"
	vFor   = "v-for"
	vIf    = "v-if"
	vModel = "v-model"
	vOn    = "v-on"
)

type renderer struct {
	tmpl []byte
	el   dom.Element
	cbs  *callbacks
	id   int64
	tree *vdom.Tree
	flag *html.Node
}

// newRenderer creates a new renderer.
func newRenderer(el string, tmpl []byte, cbs *callbacks) *renderer {
	element := dom.GetWindow().Document().QuerySelector(el)

	minifier := minify.New()
	minifier.Add("text/html", &minhtml.Minifier{KeepEndTags: true})
	tmpl, err := minifier.Bytes("text/html", tmpl)
	must(err)

	return &renderer{el: element, tmpl: tmpl, cbs: cbs, tree: &vdom.Tree{}, flag: &html.Node{}}
}

// render executes the template with the given data and applies it to the dom element.
func (renderer *renderer) render(data map[string]interface{}) {
	buf := bytes.NewBuffer(renderer.tmpl)
	nodes := parse(buf)
	if n := len(nodes); n != 1 {
		must(fmt.Errorf("template must have a single root element, found %d", n))
	}

	node := renderer.renderNode(nodes[0], data)

	buf = bytes.NewBuffer(nil)
	err := html.Render(buf, node)
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

// parse parses the template into html nodes.
func parse(reader io.Reader) []*html.Node {
	nodes, err := html.ParseFragment(reader, &html.Node{
		Type:     html.ElementNode,
		Data:     "div",
		DataAtom: atom.Div,
	})
	must(err)
	return nodes
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
	case vFor:
		node = renderer.renderAttrFor(node, attr.Val, data)
	case vBind:
		renderAttrBind(node, part, attr.Val)
	case vModel:
		renderer.renderAttrModel(node, attr.Val, data)
	case vOn:
		renderer.renderAttrOn(node, part, attr.Val)
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

// renderAttrFor renders the vue for attribute.
func (renderer *renderer) renderAttrFor(node *html.Node, value string, data map[string]interface{}) *html.Node {
	vals := strings.Split(value, "in")
	name := bytes.TrimSpace([]byte(vals[0]))
	field := strings.TrimSpace(vals[1])

	slice, ok := data[field]
	if !ok {
		must(fmt.Errorf("slice not found for field: %s", field))
	}

	elem := bytes.NewBuffer(nil)
	err := html.Render(elem, node)
	must(err)

	buf := bytes.NewBuffer(nil)
	values := reflect.ValueOf(slice)
	n := values.Len()
	for i := 0; i < n; i++ {
		key := fmt.Sprintf("%s%d", name, renderer.id)
		renderer.id++

		b := bytes.Replace(elem.Bytes(), name, []byte(key), -1)
		_, err := buf.Write(b)
		must(err)

		data[key] = values.Index(i).Interface()
	}

	nodes := parse(buf)
	node.Parent.InsertBefore(renderer.flag, node)
	for _, child := range nodes {
		node.Parent.InsertBefore(child, node)
	}
	node.Parent.RemoveChild(node)

	return renderer.flag
}

// renderAttrBind renders the vue bind attribute.
func renderAttrBind(node *html.Node, key, value string) {
	node.Attr = append(node.Attr, html.Attribute{Key: key, Val: fmt.Sprintf("{{ %v }}", value)})
}

// renderAttrModel renders the vue model attribute.
func (renderer *renderer) renderAttrModel(node *html.Node, field string, data map[string]interface{}) {
	typ := "input"
	node.Attr = append(node.Attr, html.Attribute{Key: typ, Val: field})
	renderer.cbs.addCallback(renderer.el, typ, renderer.cbs.vModel)

	value, ok := data[field]
	if !ok {
		must(fmt.Errorf("unknown data field: %s", field))
	}
	val, ok := value.(string)
	if !ok {
		must(fmt.Errorf("data field is not of type string: %T", field))
	}
	node.Attr = append(node.Attr, html.Attribute{Key: "value", Val: val})
}

// renderAttrOn renders the vue on attribute.
func (renderer *renderer) renderAttrOn(node *html.Node, typ, method string) {
	node.Attr = append(node.Attr, html.Attribute{Key: typ, Val: method})
	renderer.cbs.addCallback(renderer.el, typ, renderer.cbs.vOn)
}

// deleteAttr deletes the attribute of the node at the index.
// Attribute order is not preserved.
func deleteAttr(node *html.Node, i int) {
	n := len(node.Attr) - 1
	node.Attr[i] = node.Attr[n]
	node.Attr = node.Attr[:n]
}
