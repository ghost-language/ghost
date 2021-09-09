package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateGet(node *ast.Get, env *object.Environment) (object.Object, bool) {
	value, ok := Evaluate(node.Expression, env)

	if !ok {
		return nil, ok
	}

	switch gettable := value.(type) {
	case *object.ClassInstance:
		// fmt.Printf("Got class instance: %s -> %s\n", gettable.Class.String(), node.Name.Lexeme)
		// fmt.Printf("%v\n\n", gettable.Class.Methods)

		if method, ok := gettable.Class.Methods[node.Name.Lexeme]; ok {
			// fmt.Println("found method")
			return method, ok
		}

		// default:
		// 	fmt.Printf("GET: %s\n", gettable.String())
	}

	// if accessor, ok := value.(object.PropertyAccessor); ok {
	// 	return accessor.Get(node.Name), true
	// }

	return nil, false
}
