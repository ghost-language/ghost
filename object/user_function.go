package object

import (
	"bytes"
	"strings"

	"ghostlang.org/x/ghost/ast"
)

// USER_FUNCTION represents the object's type.
const USER_FUNCTION = "USER_FUNCTION"

type UserFunction struct {
	// Callable
	Parameters []ast.Identifier
	Body []ast.StatementNode
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

func (uf *UserFunction) Bind(instance *ClassInstance) *UserFunction {
	thisEnv := ExtendEnvironment(uf.Env)
	thisEnv.Set("this", instance)

	return &UserFunction{
		Parameters: uf.Parameters,
		Body: uf.Body,
		Defaults: uf.Defaults,
		Env: thisEnv,
	}
}

// func (uf *UserFunction) Call(arguments []Object) {
// 	env := ExtendEnvironment(uf.Env)

// 	for i, parameter := range uf.Parameters {
// 		env.Declare(parameter.Name.Lexeme, arguments[i])
// 	}

// 	// interpreter.Evaluate(uf.Body, env)
// }

func (uf *UserFunction) Arity() int {
	return len(uf.Parameters)
}