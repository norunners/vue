package vue

// Context is received by methods to interact with the component.
type Context interface {
	Data() interface{}
	Get(field string) interface{}
	Set(field string, value interface{})
	Call(method string)
	render()
}
