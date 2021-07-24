package object

import "fmt"

// CLASS_INSTANCE represents the object's type.
const CLASS_INSTANCE = "CLASS_INSTANCE"

type ClassInstance struct {
	Class *Class
	fields map[string]interface{}
}

func (ci *ClassInstance) String() string {
	return fmt.Sprintf("<class-instance %s>", ci.Class.Name)
}

// Type returns the class instance object type.
func (ci *ClassInstance) Type() Type {
	return CLASS_INSTANCE
}