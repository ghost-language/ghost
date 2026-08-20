package object

// Scope objects consist of an environment and parent object.
type Scope struct {
	Environment *Environment
	Self        Object

	// Class records which class declared the code currently executing. It is
	// what `super` resolves against, so an inherited method always starts its
	// super lookup from its own declaring class rather than from the receiver's
	// class — otherwise a method inherited two levels down would resolve
	// `super` back to itself.
	Class *Class
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
