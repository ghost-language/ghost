package object

import "fmt"

const INSTANCE = "INSTANCE"

// Instance objects consist of a body and an environment.
type Instance struct {
	Class *Class
}

// String represents the instance object's value as a string.
func (instance *Instance) String() string {
	return fmt.Sprintf("class instance %s", instance.Class.Name.Value)
}

// Type returns the instance object type.
func (instance *Instance) Type() Type {
	return INSTANCE
}

// Method defines the set of methods available on instance objects.
func (instance *Instance) Method(method string, args []Object) (Object, bool) {
	return nil, false
}
