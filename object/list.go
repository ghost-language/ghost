package object

import (
	"bytes"
	"sort"
	"strings"

	"ghostlang.org/x/ghost/fault"
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
	case "concat":
		return list.concat(tok, args)
	case "contains":
		return list.contains(tok, args)
	case "each":
		return list.each(tok, args)
	case "filter":
		return list.filter(tok, args)
	case "first":
		return list.first(tok, args)
	case "join":
		return list.join(tok, args)
	case "last":
		return list.last(tok, args)
	case "length":
		return list.length(tok, args)
	case "map":
		return list.mapElements(tok, args)
	case "pop":
		return list.pop(tok, args)
	case "push":
		return list.push(tok, args)
	case "reduce":
		return list.reduce(tok, args)
	case "reverse":
		return list.reverse(tok, args)
	case "shift":
		return list.shift(tok, args)
	case "slice":
		return list.slice(tok, args)
	case "sort":
		return list.sort(tok, args)
	case "tail":
		return list.tail(tok, args)
	case "toString":
		return list.toString(tok, args)
	case "unique":
		return list.unique(tok, args)
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

// contains reports whether a value appears in the list, comparing contents
// rather than identity - the same rule `==` follows.
func (list *List) contains(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("list.contains()", tok, args, 1); err != nil {
		return err, true
	}

	needle, err := AnyArgument("list.contains()", tok, args, 0)

	if err != nil {
		return err, true
	}

	for _, element := range list.Elements {
		if ValuesEqual(element, needle) {
			return &Boolean{Value: true}, true
		}
	}

	return &Boolean{Value: false}, true
}

// each calls a function once per element, for the side effect rather than the
// result, and answers with the list itself so a call can chain.
func (list *List) each(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("list.each()", tok, args, 1); err != nil {
		return err, true
	}

	fn, err := FunctionArgument("list.each()", tok, args, 0)

	if err != nil {
		return err, true
	}

	for index, element := range list.Elements {
		result := fn.Call([]Object{element, NewInt(int64(index))})

		if IsError(result) {
			return result, true
		}
	}

	return list, true
}

// filter keeps the elements a function answers true for, in the same order,
// and answers with a new list.
func (list *List) filter(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("list.filter()", tok, args, 1); err != nil {
		return err, true
	}

	fn, err := FunctionArgument("list.filter()", tok, args, 0)

	if err != nil {
		return err, true
	}

	elements := make([]Object, 0, len(list.Elements))

	for index, element := range list.Elements {
		result := fn.Call([]Object{element, NewInt(int64(index))})

		if IsError(result) {
			return result, true
		}

		if isTruthy(result) {
			elements = append(elements, element)
		}
	}

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

// mapElements answers a new list built by calling a function on every
// element, in order. It is named mapElements on the Go side only because
// `map` is a keyword; the language sees it as `map`.
func (list *List) mapElements(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("list.map()", tok, args, 1); err != nil {
		return err, true
	}

	fn, err := FunctionArgument("list.map()", tok, args, 0)

	if err != nil {
		return err, true
	}

	elements := make([]Object, len(list.Elements))

	for index, element := range list.Elements {
		result := fn.Call([]Object{element, NewInt(int64(index))})

		if IsError(result) {
			return result, true
		}

		elements[index] = result
	}

	return &List{Elements: elements}, true
}

// pop removes and answers the last element, mirroring push. A list used as a
// stack grows and shrinks from the same end.
func (list *List) pop(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("list.pop()", tok, args, 0); err != nil {
		return err, true
	}

	length := len(list.Elements)

	if length == 0 {
		return &Null{}, true
	}

	last := list.Elements[length-1]
	list.Elements = list.Elements[:length-1]

	return last, true
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

// reduce folds the list down to a single value, left to right. Without a
// starting value the first element plays that role, which is why an empty
// list needs one supplied.
func (list *List) reduce(tok token.Token, args []Object) (Object, bool) {
	if err := ArityRange("list.reduce()", tok, args, 1, 2); err != nil {
		return err, true
	}

	fn, err := FunctionArgument("list.reduce()", tok, args, 0)

	if err != nil {
		return err, true
	}

	elements := list.Elements
	var accumulator Object

	if len(args) == 2 {
		accumulator = args[1]
	} else {
		if len(elements) == 0 {
			return NewError(fault.Argument, tok, "`list.reduce()` needs an initial value to reduce an empty list"), true
		}

		accumulator = elements[0]
		elements = elements[1:]
	}

	for index, element := range elements {
		accumulator = fn.Call([]Object{accumulator, element, NewInt(int64(index))})

		if IsError(accumulator) {
			return accumulator, true
		}
	}

	return accumulator, true
}

func (list *List) reverse(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("list.reverse()", tok, args, 0); err != nil {
		return err, true
	}

	length := len(list.Elements)
	elements := make([]Object, length)

	for index, element := range list.Elements {
		elements[length-1-index] = element
	}

	return &List{Elements: elements}, true
}

// shift removes and answers the first element, sliding everything else down.
// It is what pop() did before pop() was made to match push() and act on the
// other end of the list.
func (list *List) shift(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("list.shift()", tok, args, 0); err != nil {
		return err, true
	}

	if len(list.Elements) == 0 {
		return &Null{}, true
	}

	first := list.Elements[0]
	list.Elements = list.Elements[1:]

	return first, true
}

// slice answers a new list holding the elements from start up to, but not
// including, end - which defaults to the length of the list.
func (list *List) slice(tok token.Token, args []Object) (Object, bool) {
	if err := ArityRange("list.slice()", tok, args, 1, 2); err != nil {
		return err, true
	}

	start, err := NumberArgument("list.slice()", tok, args, 0)

	if err != nil {
		return err, true
	}

	length := int64(len(list.Elements))
	from := start.Int64()
	to := length

	if len(args) == 2 {
		end, err := NumberArgument("list.slice()", tok, args, 1)

		if err != nil {
			return err, true
		}

		to = end.Int64()
	}

	if from < 0 || from > length {
		return NewError(fault.Index, tok, "`list.slice()` start index %d is out of range for a list of length %d", from, length), true
	}

	if to < from || to > length {
		return NewError(fault.Index, tok, "`list.slice()` end index %d is out of range for a list of length %d", to, length), true
	}

	elements := make([]Object, to-from)
	copy(elements, list.Elements[from:to])

	return &List{Elements: elements}, true
}

// sort answers a new, sorted list, leaving this one untouched. With no
// comparator it sorts a list of numbers or of strings by their natural order;
// anything else needs a comparator, since there is no order to guess at
// otherwise. The comparator takes two elements and answers a negative, zero,
// or positive number, the same contract a comparator has everywhere else.
func (list *List) sort(tok token.Token, args []Object) (Object, bool) {
	if err := ArityRange("list.sort()", tok, args, 0, 1); err != nil {
		return err, true
	}

	if len(args) == 0 {
		return list.sortDefault(tok)
	}

	fn, err := FunctionArgument("list.sort()", tok, args, 0)

	if err != nil {
		return err, true
	}

	elements := make([]Object, len(list.Elements))
	copy(elements, list.Elements)

	var callErr Object

	sort.SliceStable(elements, func(i, j int) bool {
		if callErr != nil {
			return false
		}

		result := fn.Call([]Object{elements[i], elements[j]})

		if IsError(result) {
			callErr = result

			return false
		}

		number, ok := result.(*Number)

		if !ok {
			callErr = NewError(fault.Type, tok, "`list.sort()` comparator has to return a number, got %s", TypeName(result))

			return false
		}

		return number.IsNeg()
	})

	if callErr != nil {
		return callErr, true
	}

	return &List{Elements: elements}, true
}

func (list *List) sortDefault(tok token.Token) (Object, bool) {
	elements := make([]Object, len(list.Elements))
	copy(elements, list.Elements)

	if len(elements) < 2 {
		return &List{Elements: elements}, true
	}

	kind := elements[0].Type()

	if kind != NUMBER && kind != STRING {
		return sortNeedsComparator(tok), true
	}

	for _, element := range elements {
		if element.Type() != kind {
			return sortNeedsComparator(tok), true
		}
	}

	if kind == NUMBER {
		sort.SliceStable(elements, func(i, j int) bool {
			return elements[i].(*Number).LessThan(elements[j].(*Number))
		})
	} else {
		sort.SliceStable(elements, func(i, j int) bool {
			return elements[i].(*String).Value < elements[j].(*String).Value
		})
	}

	return &List{Elements: elements}, true
}

func sortNeedsComparator(tok token.Token) *Error {
	return NewError(fault.Argument, tok, "`list.sort()` needs a comparator to sort anything but a list of only numbers or only strings").
		WithHelp("pass a function that takes two elements and returns a negative, zero, or positive number")
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

// unique answers a new list with repeats dropped, keeping the order values
// were first seen in and comparing contents the same way contains() does.
func (list *List) unique(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("list.unique()", tok, args, 0); err != nil {
		return err, true
	}

	elements := make([]Object, 0, len(list.Elements))

	for _, element := range list.Elements {
		seen := false

		for _, kept := range elements {
			if ValuesEqual(kept, element) {
				seen = true

				break
			}
		}

		if !seen {
			elements = append(elements, element)
		}
	}

	return &List{Elements: elements}, true
}
