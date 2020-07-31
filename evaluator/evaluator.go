package evaluator

import (
	"fmt"
	"io/ioutil"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/builtins"
	"ghostlang.org/x/ghost/decimal"
	"ghostlang.org/x/ghost/lexer"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/parser"
	"ghostlang.org/x/ghost/utilities"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// Eval evaluates the node and returns an object
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.ReturnStatement:
		value := Eval(node.ReturnValue, env)

		if isError(value) {
			return value
		}

		return &object.ReturnValue{Value: value}
	case *ast.AssignmentStatement:
		identifier := evalIdentifier(node.Name, env)

		if isError(identifier) {
			return identifier
		}

		value := Eval(node.Value, env)

		if isError(value) {
			return value
		}

		object, ok := identifier.(object.Mutable)

		if !ok {
			return newError("cannot assign to %s", identifier.Type())
		}

		object.Set(value)

		return value

	// Expressions
	case *ast.BindExpression:
		value := Eval(node.Value, env)

		if isError(value) {
			return value
		}

		if identifier, ok := node.Left.(*ast.Identifier); ok {
			env.Set(identifier.Value, value)
		}

		return NULL
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.NumberLiteral:
		return &object.Number{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.MapLiteral:
		return evalMapLiteral(node, env)
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)

		if isError(right) {
			return right
		}

		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)

		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)

		if isError(right) {
			return right
		}

		return evalInfixExpression(node, node.Operator, left, right, env)
	case *ast.PostfixExpression:
		return evalPostfixExpression(node, node.Operator, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.ListLiteral:
		elements := evalExpressions(node.Elements, env)

		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}

		return &object.List{Elements: elements}
	case *ast.IndexExpression:
		left := Eval(node.Left, env)

		if isError(left) {
			return left
		}

		index := Eval(node.Index, env)

		if isError(index) {
			return index
		}

		return evalIndexExpression(left, index)
	case *ast.FunctionLiteral:
		parameters := node.Parameters
		body := node.Body
		defaults := node.Defaults
		name := "__function_" + node.Name
		function := &object.Function{Parameters: parameters, Env: env, Body: body, Defaults: defaults}

		if node.Name != "" {
			env.Set(name, function)
		}

		return function
	case *ast.WhileExpression:
		return evalWhileExpression(node, env)
	case *ast.ImportExpression:
		return evalImportExpression(node, env)
	case *ast.CallExpression:
		function := Eval(node.Callable, env)

		if isError(function) {
			return function
		}

		arguments := evalExpressions(node.Arguments, env)

		if len(arguments) == 1 && isError(arguments[0]) {
			return arguments[0]
		}

		return applyFunction(function, env, arguments)
	}

	return nil
}

func EvalModule(name string) object.Object {
	filename := utilities.FindModule(name)

	if filename == "" {
		return newError("Import Error: no module named '%s' found", name)
	}

	b, err := ioutil.ReadFile(filename)

	if err != nil {
		return newError("IO Error: error reading module '%s': %s", name, err)
	}

	l := lexer.New(string(b))
	p := parser.New(l)

	module := p.ParseProgram()

	if len(p.Errors()) != 0 {
		return newError("Parse Error: %s", p.Errors())
	}

	env := object.NewEnvironment()
	Eval(module, env)

	return env.Exported()
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}

	return false
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalExpressions(expressions []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, expression := range expressions {
		evaluated := Eval(expression, env)

		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.NUMBER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Number).Value.Neg()

	return &object.Number{Value: value}
}

func evalInfixExpression(node *ast.InfixExpression, operator string, left object.Object, right object.Object, env *object.Environment) object.Object {
	switch {
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		return evalBooleanInfixExpression(operator, left, right)
	case left.Type() == object.NUMBER_OBJ && right.Type() == object.NUMBER_OBJ:
		return evalNumberInfixExpression(node, operator, left, right, env)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalPostfixExpression(node *ast.PostfixExpression, operator string, env *object.Environment) object.Object {
	switch operator {
	case "++":
		value, ok := env.Get(node.Token.Literal)

		if !ok {
			return newError("Token literal %s is unknown", node.Token.Literal)
		}

		expression, ok := value.(*object.Number)

		if !ok {
			return newError("Invalid left-hand side expression in postfix operation")
		}

		one := decimal.NewFromInt(1)
		decimal := &object.Number{Value: expression.Value.Add(one)}
		env.Set(node.Token.Literal, decimal)

		return decimal
	case "--":
		value, ok := env.Get(node.Token.Literal)

		if !ok {
			return newError("Token literal %s is unknown", node.Token.Literal)
		}

		expression, ok := value.(*object.Number)

		if !ok {
			return newError("Invalid left-hand side expression in postfix operation")
		}

		one := decimal.NewFromInt(1)
		decimal := &object.Number{Value: expression.Value.Sub(one)}
		env.Set(node.Token.Literal, decimal)

		return decimal
	default:
		return newError("unknown operator: %s", operator)
	}
}

func evalBooleanInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftValue := left.(*object.Boolean).Value
	rightValue := right.(*object.Boolean).Value

	switch operator {
	case "and":
		return nativeBoolToBooleanObject(leftValue && rightValue)
	case "or":
		return nativeBoolToBooleanObject(leftValue || rightValue)
	case "==":
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToBooleanObject(leftValue != rightValue)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalNumberInfixExpression(node *ast.InfixExpression, operator string, left object.Object, right object.Object, env *object.Environment) object.Object {
	leftValue := left.(*object.Number).Value
	rightValue := right.(*object.Number).Value

	switch operator {
	case "+":
		return &object.Number{Value: leftValue.Add(rightValue)}
	case "-":
		return &object.Number{Value: leftValue.Sub(rightValue)}
	case "*":
		return &object.Number{Value: leftValue.Mul(rightValue)}
	case "/":
		return &object.Number{Value: leftValue.Div(rightValue)}
	case "%":
		return &object.Number{Value: leftValue.Mod(rightValue)}
	case "<":
		return nativeBoolToBooleanObject(leftValue.LessThan(rightValue))
	case ">":
		return nativeBoolToBooleanObject(leftValue.GreaterThan(rightValue))
	case "<=":
		return nativeBoolToBooleanObject(leftValue.LessThanOrEqual(rightValue))
	case ">=":
		return nativeBoolToBooleanObject(leftValue.GreaterThanOrEqual(rightValue))
	case "==":
		return nativeBoolToBooleanObject(leftValue.Equal(rightValue))
	case "!=":
		return nativeBoolToBooleanObject(!leftValue.Equal(rightValue))
	case "+=":
		_, ok := env.Get(node.Left.String())

		if !ok {
			return newError("Variable %s is unknown", node.Left.String())
		}

		decimal := &object.Number{Value: leftValue.Add(rightValue)}
		env.Set(node.Left.String(), decimal)

		return NULL
	case "-=":
		_, ok := env.Get(node.Left.String())

		if !ok {
			return newError("Variable %s is unknown", node.Left.String())
		}

		decimal := &object.Number{Value: leftValue.Sub(rightValue)}
		env.Set(node.Left.String(), decimal)

		return NULL
	case "*=":
		_, ok := env.Get(node.Left.String())

		if !ok {
			return newError("Variable %s is unknown", node.Left.String())
		}

		decimal := &object.Number{Value: leftValue.Mul(rightValue)}
		env.Set(node.Left.String(), decimal)

		return NULL
	case "/=":
		_, ok := env.Get(node.Left.String())

		if !ok {
			return newError("Variable %s is unknown", node.Left.String())
		}

		decimal := &object.Number{Value: leftValue.Div(rightValue)}
		env.Set(node.Left.String(), decimal)

		return NULL
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftValue + rightValue}
	case "==":
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToBooleanObject(leftValue != rightValue)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if value, ok := env.Get(node.Value); ok {
		return value
	}

	if builtin, ok := builtins.Builtins[node.Value]; ok {
		return builtin
	}

	if function, ok := env.Get("__function_" + node.Value); ok {
		return function
	}

	return newError("identifier not found: " + node.Value)
}

func evalIndexExpression(left object.Object, index object.Object) object.Object {
	switch {
	case left.Type() == object.LIST_OBJ && index.Type() == object.NUMBER_OBJ:
		return evalListIndexExpression(left, index)
	case left.Type() == object.MAP_OBJ:
		return evalMapIndexExpression(left, index)
	case left.Type() == object.MODULE_OBJ:
		return evalModuleIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalListIndexExpression(list object.Object, index object.Object) object.Object {
	listObject := list.(*object.List)
	idx := index.(*object.Number).Value.IntPart()
	max := int64(len(listObject.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return listObject.Elements[idx]
}

func evalMapLiteral(node *ast.MapLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.MapKey]object.MapPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)

		if isError(key) {
			return key
		}

		mapKey, ok := key.(object.Mappable)

		if !ok {
			return newError("unusable as map key: %s", key.Type())
		}

		value := Eval(valueNode, env)

		if isError(value) {
			return value
		}

		mapped := mapKey.MapKey()
		pairs[mapped] = object.MapPair{Key: key, Value: value}
	}

	return &object.Map{Pairs: pairs}
}

func evalMapIndexExpression(m object.Object, index object.Object) object.Object {
	mapObject := m.(*object.Map)

	key, ok := index.(object.Mappable)

	if !ok {
		return newError("unusable as map key: %s", index.Type())
	}

	pair, ok := mapObject.Pairs[key.MapKey()]

	if !ok {
		return NULL
	}

	return pair.Value
}

func evalModuleIndexExpression(module, index object.Object) object.Object {
	moduleObject := module.(*object.Module)

	return evalMapIndexExpression(moduleObject.Attributes, index)
}

func evalWhileExpression(we *ast.WhileExpression, env *object.Environment) object.Object {
	for {
		condition := Eval(we.Condition, env)

		if isError(condition) {
			return condition
		}

		if isTruthy(condition) {
			Eval(we.Consequence, env)
		} else {
			break
		}
	}

	return NULL
}

func evalImportExpression(ie *ast.ImportExpression, env *object.Environment) object.Object {
	name := Eval(ie.Name, env)

	if isError(name) {
		return name
	}

	if s, ok := name.(*object.String); ok {
		attributes := EvalModule(s.Value)

		if isError(attributes) {
			return attributes
		}

		return &object.Module{Name: s.Value, Attributes: attributes}
	}

	return newError("Import Error: invalid import path '%s'", name)
}

func applyFunction(fn object.Object, env *object.Environment, arguments []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, arguments)
		evaluated := Eval(fn.Body, extendedEnv)

		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		if result := fn.Fn(env, arguments...); result != nil {
			return result
		}

		return NULL
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, arguments []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for key, value := range fn.Defaults {
		env.Set(key, Eval(value, env))
	}

	for index, parameter := range fn.Parameters {
		if index < len(arguments) {
			env.Set(parameter.Value, arguments[index])
		}
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}

	return FALSE
}
