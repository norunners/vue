package vue

// Methods is the vue methods option for components.
// See API: https://vuejs.org/v2/api/#methods for details.
func Methods(options ...Option) Option {
	methods := Make(options...)
	return Field("methods", methods)
}
