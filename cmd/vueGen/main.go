// Package main is vueGen, the go code generator for vue templates.
// Files with the .vue extension are parsed and used to generate go template files
// where the template body is the value of an unexported constant
// with the name as the filebase of the vue template file.
package main

import (
	"encoding/xml"
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// VueGen contains the decoded xml of the template file.
type VueGen struct {
	Template Template `xml:"template"`
	Style    Style    `xml:"style"`
}

// Template contains the inner xml of the template element.
type Template struct {
	Value string `xml:",innerxml"`
}

// Style is not currently supported.
type Style struct {
	Value  string `xml:",innerxml"`
	Scoped bool   `xml:"scoped,attr"`
}

func main() {
	gofile := os.Getenv("GOFILE")
	gopackage := os.Getenv("GOPACKAGE")

	minifier := minify.New()
	minifier.Add("text/html", &html.Minifier{
		KeepConditionalComments: true,
		KeepDefaultAttrVals:     true,
		KeepDocumentTags:        true,
		KeepEndTags:             true,
	})

	templates, err := filepath.Glob("*.vue")
	must(err)
	for _, template := range templates {
		filebase := strings.TrimSuffix(template, ".vue")
		filename := filebase + ".go"
		// Prevent overwriting the client of the generated source file.
		if filename == gofile {
			panic(fmt.Errorf("file conflict on name: %s", filename))
		}

		// Read the vue template file.
		bytes, err := ioutil.ReadFile(template)
		must(err)

		// Decode the xml template.
		data := fmt.Sprintf("<vueGen>%s</vueGen>", string(bytes))
		dec := xml.NewDecoder(strings.NewReader(data))
		dec.Strict = false
		vueGen := &VueGen{}
		err = dec.Decode(vueGen)
		must(err)

		// Minify the template to remove extraneous whitespace.
		value, err := minifier.String("text/html", vueGen.Template.Value)
		must(err)

		// Generate the source code.
		source := jen.NewFile(gopackage)
		comment := fmt.Sprintf("This source was generated with vueGen from file: %s, do not edit.", template)
		source.HeaderComment(comment)
		source.Line()
		// Ensure the constant name is unexported.
		name := strings.ToLower(filebase[:1]) + filebase[1:]
		source.Const().Id(name).Op("=").Lit(value)

		// Write the source to the go file.
		file, err := os.Create(filename)
		must(err)
		err = source.Render(file)
		must(err)
		err = file.Close()
		must(err)
	}
}

// must panics on errors.
func must(err error) {
	if err != nil {
		panic(err)
	}
}
