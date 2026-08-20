package object

// Super is the receiver produced by the `super` keyword. It carries the
// instance the call is bound to along with the class member lookup should start
// from, which is the superclass of the class that declared the running method.
type Super struct {
	Instance *Instance
	Class    *Class
}

// String represents the super object's value as a string.
func (super *Super) String() string {
	return "super"
}

// Type returns the super object type.
func (super *Super) Type() Type {
	return SUPER
}

// Method defines the set of methods available on super objects.
func (super *Super) Method(method string, args []Object) (Object, bool) {
	return nil, false
}
