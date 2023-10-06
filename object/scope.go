package object

const SCOPE = "SCOPE"

// Scope objects consist of an environment and parent object.
type Scope struct {
	Environment *Environment
	Self        Object
}

// String represents the scope object's value as a string.
func (scope *Scope) String() string {
	return "scope"
}

// Type returns the scope object type.
func (scope *Scope) Type() Type {
	return SCOPE
}

// Method defines the set of methods available on scope objects.
func (scope *Scope) Method(method string, args []Object) (Object, bool) {
	return nil, false
}
