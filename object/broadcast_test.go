package object

import "testing"

func numbers(values ...int64) *List {
	elements := make([]Object, len(values))

	for index, value := range values {
		elements[index] = NewInt(value)
	}

	return &List{Elements: elements}
}

func rows(lists ...*List) *List {
	elements := make([]Object, len(lists))

	for index, list := range lists {
		elements[index] = list
	}

	return &List{Elements: elements}
}

// add is the simplest operation to broadcast, so the tests below are about the
// shapes rather than the arithmetic.
func add(values []*Number) Object {
	return values[0].Add(values[1])
}

func TestBroadcastShapes(t *testing.T) {
	tests := []struct {
		name     string
		left     Object
		right    Object
		expected string
	}{
		{"number and number", NewInt(2), NewInt(3), "5"},
		{"number spreads across a list", numbers(1, 2, 3), NewInt(10), "[11, 12, 13]"},
		{"either way round", NewInt(10), numbers(1, 2, 3), "[11, 12, 13]"},
		{"lists pair off", numbers(1, 2, 3), numbers(10, 20, 30), "[11, 22, 33]"},
		{"matrices pair off", rows(numbers(1, 2), numbers(3, 4)), rows(numbers(10, 20), numbers(30, 40)), "[[11, 22], [33, 44]]"},
		{"number spreads across a matrix", rows(numbers(1, 2), numbers(3, 4)), NewInt(10), "[[11, 12], [13, 14]]"},

		// The cases that make this numpy's rule rather than a positional
		// pairing: a shorter shape is aligned from the right and stretched.
		{"a row stretches down a matrix", rows(numbers(1, 2), numbers(3, 4)), numbers(10, 20), "[[11, 22], [13, 24]]"},
		{"a column stretches across a matrix", rows(numbers(1, 2), numbers(3, 4)), rows(numbers(10), numbers(20)), "[[11, 12], [23, 24]]"},
		{"an axis of one repeats", numbers(1, 2, 3), numbers(5), "[6, 7, 8]"},
		{"a row against a column spans both", rows(numbers(1), numbers(2)), numbers(10, 20), "[[11, 21], [12, 22]]"},

		{"empty lists have a shape", &List{Elements: []Object{}}, &List{Elements: []Object{}}, "[]"},
	}

	for _, test := range tests {
		result, fault := Broadcast([]Object{test.left, test.right}, add)

		if fault != nil {
			t.Errorf("%s: unexpected fault: %s", test.name, fault.Reason)
			continue
		}

		if result.String() != test.expected {
			t.Errorf("%s: got=%s, expected=%s", test.name, result.String(), test.expected)
		}
	}
}

func TestBroadcastFaults(t *testing.T) {
	tests := []struct {
		name     string
		left     Object
		right    Object
		expected string
	}{
		{
			"shapes that cannot be lined up",
			numbers(1, 2),
			numbers(1, 2, 3),
			"shapes 2 and 3 cannot be combined",
		},
		{
			"matrices of different widths",
			rows(numbers(1, 2), numbers(3, 4)),
			rows(numbers(1, 2, 3), numbers(4, 5, 6)),
			"shapes 2×2 and 2×3 cannot be combined",
		},
		{
			"ragged lists have no shape",
			rows(numbers(1, 2), numbers(3)),
			NewInt(1),
			"lists have to be rectangular to combine elementwise",
		},
		{
			"elements that are not numbers",
			&List{Elements: []Object{&String{Value: "a"}}},
			NewInt(1),
			"expected a number or a list of numbers, found STRING",
		},
	}

	for _, test := range tests {
		_, fault := Broadcast([]Object{test.left, test.right}, add)

		if fault == nil {
			t.Errorf("%s: expected a fault", test.name)
			continue
		}

		if fault.Reason != test.expected {
			t.Errorf("%s: got=%s, expected=%s", test.name, fault.Reason, test.expected)
		}
	}
}

// TestBroadcastAcrossThreeOperands covers the arity the three-argument methods
// need, where clamp holds two bounds against a whole matrix.
func TestBroadcastAcrossThreeOperands(t *testing.T) {
	middle := func(values []*Number) Object {
		if values[0].LessThan(values[1]) {
			return values[1]
		}

		if values[0].GreaterThan(values[2]) {
			return values[2]
		}

		return values[0]
	}

	result, fault := Broadcast([]Object{rows(numbers(-5, 5), numbers(15, 3)), NewInt(0), NewInt(10)}, middle)

	if fault != nil {
		t.Fatalf("unexpected fault: %s", fault.Reason)
	}

	if result.String() != "[[0, 5], [10, 3]]" {
		t.Errorf("got=%s, expected=[[0, 5], [10, 3]]", result.String())
	}
}

// TestBroadcastStopsAtErrors covers an operation that fails partway, such as a
// division by zero inside a list.
func TestBroadcastStopsAtErrors(t *testing.T) {
	divide := func(values []*Number) Object {
		if values[1].IsZero() {
			return NewError("division by zero")
		}

		return values[0].Div(values[1])
	}

	result, fault := Broadcast([]Object{numbers(1, 2, 3), numbers(1, 0, 3)}, divide)

	if fault != nil {
		t.Fatalf("unexpected fault: %s", fault.Reason)
	}

	if !IsError(result) {
		t.Fatalf("expected an error, got=%s", result.String())
	}
}
