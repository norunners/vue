vue
===
[![GoDoc](https://godoc.org/github.com/norunners/vue?status.svg)](https://godoc.org/github.com/norunners/vue)
[![Build Status](https://travis-ci.org/norunners/vue.svg?branch=master)](https://travis-ci.org/norunners/vue)

Package `vue` is the progressive framework for [gopherjs](https://github.com/gopherjs/gopherjs) applications.

Install
-------
```bash
go get -u github.com/norunners/vue
```

Goals
-----
* Provide a cohesive solution for a framework, state manager and front end router.
* Encourage component reuse to promote organization and combat monolithic front end applications.
* Leverage [templating](https://github.com/norunners/vue/blob/master/cmd/vueGen/README.md) to enable developers to focus on application logic with gopherjs and not front end rendering.
* Add type safety where possible while keeping the gopherjs binding very close to the vue API.

Hello World!
------------
The `main.go` source file.
```go
type data struct {
	*js.Object

	message string `js:"message"`
}

func main() {
	data := &data{Object: js.Global.Get("Object").New()}
	data.message = "Hello World!"

	app := vue.New(
		vue.El("#app"),
		vue.DataValue(data),
	)
	js.Global.Set("app", app)
}
```

The `index.html` templating file.
```html
<!DOCTYPE html>
<html>
    <head>
        <meta charset="UTF-8">
        <title>Hello World!</title>
        <script type="text/javascript" src="path/to/vue.js"></script>
    </head>
    <body>
        <div id="app">
            {{ message }}
        </div>
        <script type="text/javascript" src="path/to/main.js"></script>
    </body>
</html>
```

Serve Examples
------------------
Run the following command and navigate to: [localhost:8080/examples](http://localhost:8080/examples)
```bash
gopherjs serve github.com/norunners/vue
```
Then navigate to an example, e.g [Hello World!](http://localhost:8080/examples/a-introduction/a-declarative-rendering)

#### Exercise the Data Binding
Open the js console in the browser and set the message to a new value.
```js
app.message = "Hello gopherjs vue!"
```

Status
------
Alpha - The state of this project is experimental until the vue API surface area is scoped out.
Therefore, adding new features of the vue API is a priority while ensuring minimal API breaking changes.

F.A.Q.
------

#### Why was vue chosen?
One of the common themes of existing frameworks is to prefer rendering HTML and components within the gopherjs source.
Vue by nature renders HTML and components as templates in a progressive manor.
Therefore, vue has a nice balance with gopherjs as the HTML and components are written within templating files and application logic remains within the gopherjs source.

#### Why build a new gopherjs binding for vue?
First and foremost, there is no reason to critique any existing vue bindings and currently do not have an exhaustive list of differences between the existing projects.
However, the existing vue bindings to not appear to share the goals of this project, which are stated above.
Final note, everyone should be encouraged to experiment with the existing gopherjs frameworks and vue bindings and choose whichever one suites their project best.

License
-------
* [MIT License](LICENSE)
