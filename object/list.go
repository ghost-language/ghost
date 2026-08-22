package object

import (
	"bytes"
	"strings"

	"ghostlang.org/x/ghost/token"
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
func (list *List) Method(method string, tok token.Token, args []Object) (Object, bool) {
	switch method {
	case "first":
		return list.first(tok, args)
	case "concat":
		return list.concat(tok, args)
	case "join":
		return list.join(tok, args)
	case "last":
		return list.last(tok, args)
	case "length":
		return list.length(tok, args)
	case "pop":
		return list.pop(tok, args)
	case "push":
		return list.push(tok, args)
	case "tail":
		return list.tail(tok, args)
	case "toString":
		return list.toString(tok, args)
	}

	return nil, false
}

// =============================================================================
// Object methods

func (list *List) first(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("list.first()", tok, args, 0); err != nil {
		return err, true
	}

	if len(list.Elements) == 0 {
		return &Null{}, true
	}

	return list.Elements[0], true
}

// concat joins two lists end to end, answering with a new list and leaving both
// operands alone. Joining is what `+` means for strings but not for lists,
// where the operators are elementwise arithmetic instead.
func (list *List) concat(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("list.concat()", tok, args, 1); err != nil {
		return err, true
	}

	other, err := ListArgument("list.concat()", tok, args, 0)

	if err != nil {
		return err.WithHelp("`concat` joins two lists; to repeat a list, multiply it"), true
	}

	elements := make([]Object, 0, len(list.Elements)+len(other.Elements))
	elements = append(elements, list.Elements...)
	elements = append(elements, other.Elements...)

	return &List{Elements: elements}, true
}

func (list *List) join(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("list.join()", tok, args, 1); err != nil {
		return err, true
	}

	separator, err := StringArgument("list.join()", tok, args, 0)

	if err != nil {
		return err, true
	}

	parts := make([]string, len(list.Elements))

	for index, value := range list.Elements {
		parts[index] = value.String()
	}

	return &String{Value: strings.Join(parts, separator.Value)}, true
}

func (list *List) last(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("list.last()", tok, args, 0); err != nil {
		return err, true
	}

	length := len(list.Elements)

	if length == 0 {
		return &Null{}, true
	}

	return list.Elements[length-1], true
}

func (list *List) length(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("list.length()", tok, args, 0); err != nil {
		return err, true
	}

	return NewInt(int64(len(list.Elements))), true
}

func (list *List) pop(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("list.pop()", tok, args, 0); err != nil {
		return err, true
	}

	if len(list.Elements) == 0 {
		return &Null{}, true
	}

	first := list.Elements[0]
	list.Elements = list.Elements[1:]

	return first, true
}

func (list *List) push(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("list.push()", tok, args, 1); err != nil {
		return err, true
	}

	value, err := AnyArgument("list.push()", tok, args, 0)

	if err != nil {
		return err, true
	}

	list.Elements = append(list.Elements, value)

	return NewInt(int64(len(list.Elements))), true
}

func (list *List) tail(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("list.tail()", tok, args, 0); err != nil {
		return err, true
	}

	length := len(list.Elements)

	if length == 0 {
		return &Null{}, true
	}

	elements := make([]Object, length-1)
	copy(elements, list.Elements[1:length])

	return &List{Elements: elements}, true
}

func (list *List) toString(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("list.toString()", tok, args, 0); err != nil {
		return err, true
	}

	return &String{Value: list.String()}, true
}
