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
	case *Map:
		return MapsEqual(left, right.(*Map))
	case *Duration:
		return DurationsEqual(left, right.(*Duration))
	case *Date:
		// The instant, not every field - a Date attached to two different time
		// zones can still be equal, the same reading `<`/`>`/`==` between two
		// dates already has via evaluateDateInfix (evaluator/date.go), which
		// reaches this same comparison for a direct `date1 == date2` instead of
		// going through here - two entry points computing the one answer,
		// rather than two answers.
		return left.Time.Equal(right.(*Date).Time)
	}

	// Everything else - instances, functions, classes - compares by identity,
	// the same rule `==` applies to instances outside a list.
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

// MapsEqual compares two maps by their contents rather than by identity, to
// any depth - the same reasoning as ListsEqual. Two maps are equal when they
// have the same keys, each with an equal value; MapKey already identifies
// "the same key" the same way indexing a map does, so no separate key
// comparison is needed.
func MapsEqual(left *Map, right *Map) bool {
	if len(left.Pairs) != len(right.Pairs) {
		return false
	}

	for key, pair := range left.Pairs {
		other, ok := right.Pairs[key]

		if !ok || !ValuesEqual(pair.Value, other.Value) {
			return false
		}
	}

	return true
}

// DurationsEqual compares two Durations field by field - a Duration is a
// small immutable record, not an identity, so two built separately with the
// same components should compare equal, the same reasoning as ListsEqual and
// MapsEqual. Ordering (`< > <= >=`) stays unsupported for Duration - unlike
// two instants, "which of these two spans is longer" has no single answer
// without a reference date (a month is a different number of days depending
// on which month), the same reason Temporal requires one for calendar unit
// comparisons.
func DurationsEqual(left *Duration, right *Duration) bool {
	return left.Years == right.Years &&
		left.Months == right.Months &&
		left.Days == right.Days &&
		left.Hours == right.Hours &&
		left.Minutes == right.Minutes &&
		left.Seconds == right.Seconds
}
