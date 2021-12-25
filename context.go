package vue

// Context is received by methods to interact with the component.
type Context interface {
	Data() interface{}
	Call(name string)
}
