package object

// ValuesEqual compares two values by content rather than identity, to any
// depth. It is what `==` means between two Ghost values of the same type, and
// it is the one place that comparison is written: the evaluator's `==`
// operator and List's contains()/unique() both call it, so a value counted as
// equal inside a list is equal everywhere else in the language too.
func ValuesEqual(left Object, right Object) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}

	if left.Type() != right.Type() {
		return false
	}

	switch left := left.(type) {
	case *Number:
		return left.Equal(right.(*Number))
	case *String:
		return left.Value == right.(*String).Value
	case *Boolean:
		return left.Value == right.(*Boolean).Value
	case *Null:
		return true
	case *List:
		return ListsEqual(left, right.(*List))
	}

	// Everything else - instances, functions, maps - compares by identity, the
	// same rule `==` applies to instances outside a list.
	return left == right
}

// ListsEqual compares two lists by their contents rather than by identity, to
// any depth.
func ListsEqual(left *List, right *List) bool {
	if len(left.Elements) != len(right.Elements) {
		return false
	}

	for index, element := range left.Elements {
		if !ValuesEqual(element, right.Elements[index]) {
			return false
		}
	}

	return true
}
