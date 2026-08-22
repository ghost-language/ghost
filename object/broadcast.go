package object

import "fmt"

// Broadcasting is what lets an operation written for single numbers apply to
// whole lists. It lives here, below both the evaluator and the library, because
// two callers need to agree on it exactly: `a + b` in the evaluator and
// `math.add(a, b)` in the math module are the same operation, and would drift
// apart if each carried its own rules.
//
// The rules are numpy's. Shapes are lined up from the right; at each axis the
// two lengths must match, or one of them must be 1, or one operand must have no
// axis there at all. Whichever it is, the shorter side repeats:
//
//	[1, 2, 3] + 10                    the number repeats across the list
//	[1, 2, 3] + [10, 20, 30]          paired off, axis by axis
//	[[1, 2], [3, 4]] + [10, 20]       the row repeats down the matrix
//
// The last case is the one worth knowing. The row is stretched across every row
// of the matrix rather than paired against them, which is what a reader coming
// from numpy expects and what makes a bias vector work on a batch of samples.

// BroadcastFault explains why two values could not be combined. It carries no
// position information: the caller knows where it is and what it was called,
// and formats the message around this.
type BroadcastFault struct {
	Reason string
}

func (fault *BroadcastFault) Error() string {
	return fault.Reason
}

// Broadcast applies an operation across its operands elementwise, stretching
// the smaller shapes across the larger. The operation is handed one number per
// operand and may answer with any object, including an *Error, which stops the
// walk and is returned as-is.
func Broadcast(operands []Object, operation func(values []*Number) Object) (Object, *BroadcastFault) {
	shapes := make([][]int, len(operands))
	result := []int(nil)

	for index, operand := range operands {
		shape, fault := shapeOf(operand)

		if fault != nil {
			return nil, fault
		}

		shapes[index] = shape

		combined, ok := combineShapes(result, shape)

		if !ok {
			return nil, &BroadcastFault{Reason: fmt.Sprintf("shapes %s and %s cannot be combined", describeShape(result), describeShape(shape))}
		}

		result = combined
	}

	return walk(operands, shapes, result, make([]*Number, len(operands)), operation), nil
}

// walk descends one axis of the result at a time. Every operand is carried
// along with its own remaining shape, so each one knows whether it has an axis
// here to index into, a single element to repeat, or nothing at this level at
// all and applies whole.
func walk(operands []Object, shapes [][]int, shape []int, values []*Number, operation func([]*Number) Object) Object {
	if len(shape) == 0 {
		for index, operand := range operands {
			values[index] = operand.(*Number)
		}

		return operation(values)
	}

	elements := make([]Object, shape[0])
	inner := make([]Object, len(operands))
	innerShapes := make([][]int, len(operands))

	for position := 0; position < shape[0]; position++ {
		for index, operand := range operands {
			inner[index], innerShapes[index] = elementAt(operand, shapes[index], len(shape), position)
		}

		result := walk(inner, innerShapes, shape[1:], values, operation)

		if IsError(result) {
			return result
		}

		elements[position] = result
	}

	return &List{Elements: elements}
}

// elementAt picks out an operand's contribution to one position along the
// current axis.
func elementAt(operand Object, shape []int, rank int, position int) (Object, []int) {
	// The operand has no axis this far out, so the whole of it lines up against
	// every position: this is a number against a list, or a row against a
	// matrix.
	if len(shape) < rank {
		return operand, shape
	}

	list := operand.(*List)

	// An axis of length one repeats rather than being indexed.
	if shape[0] == 1 {
		return list.Elements[0], shape[1:]
	}

	return list.Elements[position], shape[1:]
}

// shapeOf measures a value: no axes for a number, and the length of each level
// for a list. Lists have to be rectangular to have a shape at all, so this is
// also where ragged input and non-numeric elements are caught.
func shapeOf(value Object) ([]int, *BroadcastFault) {
	switch value := value.(type) {
	case *Number:
		return nil, nil
	case *List:
		if len(value.Elements) == 0 {
			return []int{0}, nil
		}

		inner, fault := shapeOf(value.Elements[0])

		if fault != nil {
			return nil, fault
		}

		for _, element := range value.Elements[1:] {
			shape, fault := shapeOf(element)

			if fault != nil {
				return nil, fault
			}

			if !sameShape(shape, inner) {
				return nil, &BroadcastFault{Reason: "lists have to be rectangular to combine elementwise"}
			}
		}

		return append([]int{len(value.Elements)}, inner...), nil
	}

	return nil, &BroadcastFault{Reason: fmt.Sprintf("expected a number or a list of numbers, found %s", typeNameOf(value))}
}

// combineShapes lines two shapes up from the right and reports the shape their
// result takes.
func combineShapes(left []int, right []int) ([]int, bool) {
	rank := len(left)

	if len(right) > rank {
		rank = len(right)
	}

	combined := make([]int, rank)

	for axis := 0; axis < rank; axis++ {
		leftAxis := axisAt(left, rank, axis)
		rightAxis := axisAt(right, rank, axis)

		switch {
		case leftAxis == rightAxis:
			combined[axis] = leftAxis
		case leftAxis == 1:
			combined[axis] = rightAxis
		case rightAxis == 1:
			combined[axis] = leftAxis
		default:
			return nil, false
		}
	}

	return combined, true
}

// axisAt reads one axis of a shape that has been right-aligned into a wider
// rank. Axes the shape does not reach are 1, which is what makes a shorter
// shape stretch rather than fail.
func axisAt(shape []int, rank int, axis int) int {
	offset := rank - len(shape)

	if axis < offset {
		return 1
	}

	return shape[axis-offset]
}

func sameShape(left []int, right []int) bool {
	if len(left) != len(right) {
		return false
	}

	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}

	return true
}

func describeShape(shape []int) string {
	if len(shape) == 0 {
		return "a number"
	}

	text := ""

	for index, axis := range shape {
		if index > 0 {
			text += "×"
		}

		text += fmt.Sprintf("%d", axis)
	}

	return text
}

func typeNameOf(value Object) string {
	if value == nil {
		return "null"
	}

	return value.Type().String()
}
