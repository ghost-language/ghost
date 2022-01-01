package object

import (
	"bytes"
	"strings"
)

const LIST = "LIST"

// List objects consist of a nil value.
type List struct {
	Elements []Object
}

func (object *List) Accept(v Visitor) {
	v.visitList(object)
}

// String represents the list object's value as a string.
func (list *List) String() string {
	var out bytes.Buffer

	elements := []string{}

	for _, element := range list.Elements {
		elements = append(elements, element.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// Type returns the list object type.
func (list *List) Type() Type {
	return LIST
}

// Method defines the set of methods available on list objects.
func (list *List) Method(method string, args []Object) (Object, bool) {
	return nil, false
}
