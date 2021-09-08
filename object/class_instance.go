package object

import (
	"fmt"

	"ghostlang.org/x/ghost/token"
)

// CLASS_INSTANCE represents the object's type.
const CLASS_INSTANCE = "CLASS_INSTANCE"

type ClassInstance struct {
	Class *Class
	Fields map[string]interface{}
}

func (ci *ClassInstance) String() string {
	return fmt.Sprintf("<class-instance %s>", ci.Class.Name)
}

// Type returns the class instance object type.
func (ci *ClassInstance) Type() Type {
	return CLASS_INSTANCE
}

func (ci *ClassInstance) Get(name token.Token) Object {
	if value, ok := ci.Class.Methods[name.Lexeme]; ok {
		return value
	}

	return &Null{}
}