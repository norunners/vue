package vue

// Components is the vue components option for components.
// See guide: https://vuejs.org/v2/guide/components.html
// and https://vuejs.org/v2/guide/components-registration.html#Local-Registration for details.
func Components(options ...Option) Option {
	components := Make(options...)
	return Field("components", components)
}
