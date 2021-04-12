package object

import (
	"bytes"
	"strings"

	"ghostlang.org/x/ghost/ast"
)

// USER_FUNCTION represents the object's type.
const USER_FUNCTION = "USER_FUNCTION"

type UserFunction struct {
	Parameters []*ast.Variable
	Body *ast.Block
	Defaults map[string]ast.ExpressionNode
	Env *Environment
}

// String represents the string form of the native function object.
func (uf *UserFunction) String() string {
	var out bytes.Buffer

	parameters := []string{}

	for _, p := range uf.Parameters {
		parameters = append(parameters, p.Name.Lexeme)
	}

	out.WriteString("function")
	out.WriteString("(")
	out.WriteString(strings.Join(parameters, ", "))
	out.WriteString(") {\n")
	out.WriteString("body\n")
	out.WriteString("}\n")

	return out.String()
}

// Type returns the native function object type.
func (uf *UserFunction) Type() Type {
	return USER_FUNCTION
}