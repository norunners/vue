vueGen
======

Command `vueGen` is the go code generator for vue templates.

Install
-------
```bash
go get -u github.com/norunners/vue/cmd/vueGen
```

Hello World!
------------
The `appTmpl.vue` vue template file.
```vue
<template>
    <div>
        {{ message }}
    </div>
</template>
```

The `main.go` usage of the generated template source.
Note, the directive `//go:generate vueGen` declares the usage of `vueGen`.
Also, the constant `appTmpl` is defined in the generated source. 
```go
//go:generate vueGen

type data struct {
	*js.Object

	message string `js:"message"`
}

func main() {
	data := &data{Object: newObject()}
	data.message = "Hello World!"

	app := vue.New(
		vue.El("#app"),
		vue.Template(appTmpl),
		vue.DataValue(data),
	)
	js.Global.Set("app", app)
}

func newObject(args ...interface{}) *js.Object {
	return js.Global.Get("Object").New(args...)
}
```

The following command executes `vueGen` to generate the go template file from the vue template file.
```bash
go generate
```

Finally, `appTmpl.go` is the generated go template file.
```go
// This source was generated with vueGen from file: appTmpl.vue, do not edit.

package main

const appTmpl = "<div>{{ message }}</div>"
````

Serve Examples
----------------------
Run the following command and navigate to: [Hello vueGen](http://localhost:8080/examples/c-templates-with-vueGen/b-hello-vueGen)
```bash
gopherjs serve github.com/norunners/vue
```
Then make changes to the vue template file, generate the go template and refresh the example page.

File Watcher
------------
The following shows the file watcher setup for `vueGen` in GoLand.
The file watcher executes `go generate` on go files after changes are made.
In turn, the go template files are generated from the current state of the vue template files.
However, this configuration does not listen to changes on the vue template files, but rather the go files.
Other IDEs that support file watchers may be configured similarly.
![file-watcher](https://user-images.githubusercontent.com/25853983/39666020-c31f6d42-5051-11e8-843d-429da849a835.png)

Styling
-------
The `<style>` block is not currently supported within vue template files, including scoped CSS (SCSS) styling.
Nonetheless, inline styling with class and style bindings are handled within the template element.
See guide: https://vuejs.org/v2/guide/class-and-style.html for details.
```vue
<template>
    <div class="static">...</div>
    <div v-bind:class="{ active: isActive }">...</div>
    <div style="color: green">...</div>
    <div v-bind:style="{ color: activeColor, fontSize: fontSize + 'px' }">...</div>
</template>
```

License
-------
* [MIT License](LICENSE)
