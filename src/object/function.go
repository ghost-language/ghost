package object

import (
	"ghostlang.org/x/ghost/ast"
)

const FUNCTION = "FUNCTION"

// Function objects consist of a user-generated function.
type Function struct {
	Parameters  []*ast.Identifier
	Body        *ast.Block
	Defaults    map[string]ast.ExpressionNode
	Environment *Environment
}

func (object *Function) Accept(v Visitor) {
	v.visitFunction(object)
}

// String represents the function object's value as a string.
func (function *Function) String() string {
	return "function"
}

// Type returns the function object type.
func (function *Function) Type() Type {
	return FUNCTION
}

func (function *Function) Method(method string, args []Object) (Object, bool) {
	return nil, false
}
