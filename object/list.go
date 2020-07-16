package object

import (
	"bytes"
	"strings"
)

type List struct {
	Elements []Object
}

func (ao *List) Type() ObjectType {
	return LIST_OBJ
}

func (ao *List) Inspect() string {
	var out bytes.Buffer

	elements := []string{}

	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

func (ao *List) Set(obj Object) {
	ao.Elements = obj.(*List).Elements
}
