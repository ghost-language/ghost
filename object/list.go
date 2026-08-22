package object

import (
	"bytes"
	"strings"
)

// List objects consist of a nil value.
type List struct {
	Elements []Object
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
	switch method {
	case "first":
		return list.first(args)
	case "concat":
		return list.concat(args)
	case "join":
		return list.join(args)
	case "last":
		return list.last(args)
	case "length":
		return list.length(args)
	case "pop":
		return list.pop(args)
	case "push":
		return list.push(args)
	case "tail":
		return list.tail(args)
	case "toString":
		return list.toString(args)
	}

	return nil, false
}

// =============================================================================
// Object methods

func (list *List) first(args []Object) (Object, bool) {
	if len(list.Elements) == 0 {
		return &Null{}, true
	}
	return list.Elements[0], true
}

// concat joins two lists end to end, answering with a new list and leaving both
// operands alone. Joining is what `+` means for strings but not for lists,
// where the operators are elementwise arithmetic instead.
func (list *List) concat(args []Object) (Object, bool) {
	if len(args) != 1 {
		return nil, false
	}

	other, ok := args[0].(*List)

	if !ok {
		return nil, false
	}

	elements := make([]Object, 0, len(list.Elements)+len(other.Elements))
	elements = append(elements, list.Elements...)
	elements = append(elements, other.Elements...)

	return &List{Elements: elements}, true
}

func (list *List) join(args []Object) (Object, bool) {
	var s []string

	for _, value := range list.Elements {
		s = append(s, value.String())
	}

	str := strings.Join(s, args[0].(*String).Value)

	return &String{Value: str}, true
}

func (list *List) last(args []Object) (Object, bool) {
	length := len(list.Elements)

	if length == 0 {
		return &Null{}, true
	}
	return list.Elements[length-1], true
}

func (list *List) length(args []Object) (Object, bool) {
	return NewInt(int64(len(list.Elements))), true
}

func (list *List) pop(args []Object) (Object, bool) {
	if len(list.Elements) > 0 {
		x := list.Elements[0]
		list.Elements = list.Elements[1:]

		return x, true
	}

	return &Null{}, true
}

func (list *List) push(args []Object) (Object, bool) {
	list.Elements = append(list.Elements, args[0])

	return NewInt(int64(len(list.Elements))), true
}

func (list *List) tail(args []Object) (Object, bool) {
	length := len(list.Elements)

	if length > 0 {
		newElements := make([]Object, length-1)
		copy(newElements, list.Elements[1:length])

		return &List{Elements: newElements}, true
	}

	return &Null{}, true
}

func (list *List) toString(args []Object) (Object, bool) {
	return &String{Value: list.String()}, true
}
