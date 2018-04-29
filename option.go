package vue

// Option is a function used to set the state of map, e.g. js object.
// The option receives a map and sets a key and value.
// Options are combined together as variadic arguments to make js objects of any structure.
type Option func(m Map)

// Field is a generic option, this allows the key and value of the field to be defined for the option.
// Most vue provided options will return a field option with a determined key and a type safe value.
// For example: El("#app") is equivalent to Field("el", "#app")
// The former is preferred as the key is determined and the value is of type string.
func Field(key string, value interface{}) Option {
	return func(m Map) {
		m[key] = value
	}
}

// El is the vue el option for vue.
// See API: https://vuejs.org/v2/api/#el for details.
func El(el string) Option {
	return Field("el", el)
}

// Template is the vue template option for components.
// See API: https://vuejs.org/v2/api/#template for details.
func Template(template string) Option {
	return Field("template", template)
}

// Data is the vue data function option for components.
// See guide: https://vuejs.org/v2/guide/components.html#data-Must-Be-a-Function for details.
func Data(data interface{}) Option {
	function := func() interface{} {
		return data
	}
	return DataValue(function)
}

// DataValue is the vue data object option for the vue instance.
// See API: https://vuejs.org/v2/api/#data for details.
func DataValue(data interface{}) Option {
	return Field("data", data)
}
