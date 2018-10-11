# vue
[![GoDoc](https://godoc.org/github.com/norunners/vue?status.svg)](https://godoc.org/github.com/norunners/vue)

Package `vue` is the progressive framework for [WebAssembly](https://github.com/golang/go/wiki/WebAssembly) applications.

## Install
```bash
go get github.com/norunners/vue
```

## Goals
* Provide a unified solution for a framework, state manager and router in the frontend space.
* Leverage [templating](https://github.com/norunners/vueg) to separate application logic from frontend rendering.
* Simplify data binding to ease the relation of state management to rendering.
* Encourage component reuse to promote development productivity.
* Follow an idiomatic Go translation of the familiar Vue API.

## Hello World!
The `main.go` file is compiled to a `.wasm` WebAssemply file.
```go
package main

import "github.com/norunners/vue"

type Data struct {
	Message string
}

func main() {
	vue.New(
		vue.El("#app"),
		vue.Template("<p>{{ Message }}</p>"),
		vue.Data(Data{Message: "Hello WebAssembly!"}),
	)

	select {}
}
```

The `index.wasmgo.html` file fetches and runs a `.wasm` WebAssemply file.
```html
<!doctype html>
<html>
    <head>
        <meta charset="utf-8">
        <script src="{{ .Script }}"></script>
    </head>
    <body>
        <div id="app"></div>
        <script src="{{ .Loader }}"></script>
    </body>
</html>
```
*Note, the example above is compatible with [wasmgo](https://github.com/dave/wasmgo).*

## Serve Examples
Install `wasmgo` to serve examples.
```bash
go get -u github.com/dave/wasmgo
```

Serve an example [locally](http://localhost:8080/).
```bash
cd examples/1-declarative-rendering
wasmgo serve
```

## Status
Alpha - The state of this project is experimental until the common features of Vue are implemented.
The plan is to follow the Vue API closely except for areas of major simplification, which may lead to a subset of the Vue API.
During this stage, the API is expected to encounter minor breaking changes but increase in stability as the project progresses.

## F.A.Q.

#### Why Vue?
One of the common themes of existing frameworks is to combine component application logic with frontend rendering.
This can lead to a confusing mental model to reason about because both concerns may be mixed together in the same logic.
By design, Vue renders components with templates which ensures application logic is developed separately from frontend rending.

Another commonality of existing frameworks is to unnecessarily expose the relation of state management to rendering in the API.
By design, Vue binds data in both directions which ensures automatic updating and rendering when state changes.

This project aims to combine the simplicity of Vue with the power of Go WebAssembly.

License
-------
* [MIT License](LICENSE)
