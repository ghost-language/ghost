package object

// RETURN represents the object's type.
const RETURN = "RETURN"

type Return struct {
	Value Object
}

// String represents the string form of the return object.
func (r *Return) String() string {
	return "return"
}

// Type returns the return object type.
func (r *Return) Type() Type {
	return RETURN
}