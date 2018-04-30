package vue

import (
	"reflect"
	"testing"
)

func TestField(t *testing.T) {
	key, value := "testKey", "testValue"
	expected := Map{key: value}

	actual := Make(Field(key, value))

	assertEqual(t, expected, actual)
}

func TestEl(t *testing.T) {
	el := "#app"
	expected := Map{"el": el}

	actual := Make(El(el))

	assertEqual(t, expected, actual)
}

func TestTemplate(t *testing.T) {
	template := "<app/>"
	expected := Map{"template": template}

	actual := Make(Template(template))

	assertEqual(t, expected, actual)
}

func TestData(t *testing.T) {
	expected := "testData"

	m := Make(Data(expected))

	// Equals does not like the function value.
	actual := m["data"].(func() interface{})()
	assertEqual(t, expected, actual)
}

func TestDataValue(t *testing.T) {
	data := "testData"
	expected := Map{"data": data}

	actual := Make(DataValue(data))

	assertEqual(t, expected, actual)
}

func assertEqual(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %#v but found: %#v", expected, actual)
	}
}
