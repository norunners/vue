package vue

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/norunners/vert"
)

// PropType is the vue props type.
type PropType string

const (
	TypeString   PropType = "String"
	TypeNumber   PropType = "Number"
	TypeBoolean  PropType = "Boolean"
	TypeFunction PropType = "Function"
	TypeObject   PropType = "Object"
	TypeArray    PropType = "Array"
	TypeSymbol   PropType = "Symbol"
)

// Props is the vue props option for components.
// See API: https://vuejs.org/v2/api/#props for details.
func Props(options ...Option) Option {
	props := Make(options...)
	return Field("props", props)
}

// PropsList is the vue props array option for components.
// Note, the given props are a variadic argument of type string.
func PropsList(props ...string) Option {
	return Field("props", props)
}

// Type is the vue type option for props.
// See guide: https://vuejs.org/v2/guide/components-props.html#Type-Checks for details.
func Type(typ PropType) Option {
	return Field("type", typ)
}

// Default is the vue default option for props.
func Default(value interface{}) Option {
	return Field("default", value)
}

// Required is the vue required option for props.
func Required() Option {
	return Field("required", true)
}

// Validator is the vue validator option for props.
// See guide: https://vuejs.org/v2/guide/components-props.html#Prop-Validation for details.
func Validator(function interface{}) Option {
	fn := vert.New(function)
	validator := func(object *js.Object) bool {
		values, err := fn.Call(object)
		must(err)
		return values[0].(bool)
	}
	return Field("validator", validator)
}
