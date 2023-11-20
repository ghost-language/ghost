package object

import "github.com/shopspring/decimal"

const NUMBER = "NUMBER"

// Number objects consist of a decimal value.
type Number struct {
	Value decimal.Decimal
}

// String represents the number object's value as a string.
func (number *Number) String() string {
	return number.Value.String()
}

// Type returns the number object type.
func (number *Number) Type() Type {
	return NUMBER
}

// MapKey defines a unique hash value for use as a map key.
func (number *Number) MapKey() MapKey {
	return MapKey{Type: number.Type(), Value: uint64(number.Value.IntPart())}
}

// Method defines the set of methods available on number objects.
func (number *Number) Method(method string, args []Object) (Object, bool) {
	switch method {
	case "round":
		return number.round(args)
	case "floor":
		return number.floor(args)
	case "toString":
		return number.toString(args)
	}

	return nil, false
}

// =============================================================================
// Object methods

func (number *Number) toString(args []Object) (Object, bool) {
	return &String{Value: number.Value.String()}, true
}

func (number *Number) round(args []Object) (Object, bool) {
	places := &Number{Value: decimal.NewFromInt(0)}

	if len(args) == 1 {
		if args[0].Type() != NUMBER {
			return nil, false
		}

		places = args[0].(*Number)
	}

	return &Number{Value: number.Value.Round(int32(places.Value.IntPart()))}, true
}

func (number *Number) floor(args []Object) (Object, bool) {
	return &Number{Value: number.Value.Floor()}, true
}
