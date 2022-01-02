package object

import (
	"bytes"
	"strings"

	"github.com/shopspring/decimal"
)

const LIST = "LIST"

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
	return list.Elements[0], true
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

	return list.Elements[length-1], true
}

func (list *List) length(args []Object) (Object, bool) {
	return &Number{Value: decimal.NewFromInt(int64(len(list.Elements)))}, true
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
	length := len(list.Elements)
	newLength := length + 1

	newElements := make([]Object, newLength)
	copy(newElements, list.Elements)
	newElements[length] = args[0]

	list.Elements = newElements

	return &Number{Value: decimal.NewFromInt(int64(newLength))}, true
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
