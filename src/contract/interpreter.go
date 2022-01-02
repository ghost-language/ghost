package contract

type AstNode interface{}
type Object interface{}
type Environment interface{}

type Evaluator func(node interface{}, env *interface{}) (interface{}, bool)

// interface {
// 	Evaluate(node Ast, env *Object) (Object, bool)
// }
