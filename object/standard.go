package object

// STANDARD represents the object's type.
const STANDARD = "STANDARD"

type Standard struct {
	Name     string
	Function StandardFunction
}

// String represents the string form of the standard function object.
func (s *Standard) String() string {
	return "standard"
}

// Type returns the standard function object type.
func (s *Standard) Type() Type {
	return STANDARD
}
