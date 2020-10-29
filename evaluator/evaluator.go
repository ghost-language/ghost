package evaluator

import (
	"io/ioutil"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/builtins"
	"ghostlang.org/x/ghost/decimal"
	"ghostlang.org/x/ghost/lexer"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/parser"
	"ghostlang.org/x/ghost/utilities"
	"ghostlang.org/x/ghost/value"
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
		val := Eval(node.ReturnValue, env)

		if utilities.IsError(val) {
			return val
		}

		return &object.ReturnValue{Value: val}
	case *ast.AssignmentStatement:
		assignment := evalAssignmentStatement(node, env)

		if utilities.IsError(assignment) {
			return assignment
		}

		return value.NULL
		// identifier := evalIdentifierLiteral(node.Name, env)

		// if utilities.IsError(identifier) {
		// 	return identifier
		// }

		// val := Eval(node.Value, env)

		// if utilities.IsError(val) {
		// 	return val
		// }

		// object, ok := identifier.(object.Mutable)

		// if !ok {
		// 	return utilities.NewError("cannot assign to %s", identifier.Type())
		// }

		// object.Set(val)

		// return val

	// Expressions
	case *ast.BindExpression:
		val := Eval(node.Value, env)

		if utilities.IsError(val) {
			return val
		}

		if identifier, ok := node.Left.(*ast.IdentifierLiteral); ok {
			env.Set(identifier.Value, val)
		}

		return value.NULL
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.NumberLiteral:
		return &object.Number{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.MapLiteral:
		return evalMapLiteral(node, env)
	case *ast.BooleanLiteral:
		return utilities.NativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)

		if utilities.IsError(right) {
			return right
		}

		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)

		if utilities.IsError(left) {
			return left
		}

		right := Eval(node.Right, env)

		if utilities.IsError(right) {
			return right
		}

		return evalInfixExpression(node, node.Operator, left, right, env)
	case *ast.PostfixExpression:
		return evalPostfixExpression(node, node.Operator, env)
	case *ast.IdentifierLiteral:
		return evalIdentifierLiteral(node, env)
	case *ast.ListLiteral:
		elements := evalExpressions(node.Elements, env)

		if len(elements) == 1 && utilities.IsError(elements[0]) {
			return elements[0]
		}

		return &object.List{Elements: elements}
	case *ast.IndexExpression:
		left := Eval(node.Left, env)

		if utilities.IsError(left) {
			return left
		}

		index := Eval(node.Index, env)

		if utilities.IsError(index) {
			return index
		}

		return evalIndexExpression(left, index)
	case *ast.FunctionLiteral:
		parameters := node.Parameters
		body := node.Body
		defaults := node.Defaults
		name := node.Name
		function := &object.Function{Parameters: parameters, Env: env, Body: body, Defaults: defaults}

		if node.Name != "" {
			env.Set(name, function)
		}

		return function
	case *ast.WhileExpression:
		return evalWhileExpression(node, env)
	case *ast.ForExpression:
		return evalForExpression(node, env)
	case *ast.ForInExpression:
		return evalForInExpression(node, env)
	case *ast.ImportExpression:
		return evalImportExpression(node, env)
	case *ast.CallExpression:
		function := Eval(node.Callable, env)

		if utilities.IsError(function) {
			return function
		}

		arguments := evalExpressions(node.Arguments, env)

		if len(arguments) == 1 && utilities.IsError(arguments[0]) {
			return arguments[0]
		}

		return applyFunction(function, env, arguments)
	}

	return nil
}

func EvalModule(name string) object.Object {
	filename := utilities.FindModule(name)

	if filename == "" {
		return utilities.NewError("Import Error: no module named '%s' found", name)
	}

	b, err := ioutil.ReadFile(filename)

	if err != nil {
		return utilities.NewError("IO Error: error reading module '%s': %s", name, err)
	}

	l := lexer.New(string(b))
	p := parser.New(l)

	module := p.ParseProgram()

	if len(p.Errors()) != 0 {
		return utilities.NewError("Parse Error: %s", p.Errors())
	}

	env := object.NewEnvironment()
	Eval(module, env)

	return env.Exported()
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

		if utilities.IsError(evaluated) {
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
		return utilities.NewError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case value.TRUE:
		return value.FALSE
	case value.FALSE:
		return value.TRUE
	case value.NULL:
		return value.TRUE
	default:
		return value.FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.NUMBER_OBJ {
		return utilities.NewError("unknown operator: -%s", right.Type())
	}

	val := right.(*object.Number).Value.Neg()

	return &object.Number{Value: val}
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
		return utilities.NativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return utilities.NativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return utilities.NewError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return utilities.NewError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalPostfixExpression(node *ast.PostfixExpression, operator string, env *object.Environment) object.Object {
	switch operator {
	case "++":
		val, ok := env.Get(node.Token.Literal)

		if !ok {
			return utilities.NewError("Token literal %s is unknown", node.Token.Literal)
		}

		expression, ok := val.(*object.Number)

		if !ok {
			return utilities.NewError("Invalid left-hand side expression in postfix operation")
		}

		one := decimal.NewFromInt(1)
		decimal := &object.Number{Value: expression.Value.Add(one)}
		env.Set(node.Token.Literal, decimal)

		return decimal
	case "--":
		val, ok := env.Get(node.Token.Literal)

		if !ok {
			return utilities.NewError("Token literal %s is unknown", node.Token.Literal)
		}

		expression, ok := val.(*object.Number)

		if !ok {
			return utilities.NewError("Invalid left-hand side expression in postfix operation")
		}

		one := decimal.NewFromInt(1)
		decimal := &object.Number{Value: expression.Value.Sub(one)}
		env.Set(node.Token.Literal, decimal)

		return decimal
	default:
		return utilities.NewError("unknown operator: %s", operator)
	}
}

func evalBooleanInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftValue := left.(*object.Boolean).Value
	rightValue := right.(*object.Boolean).Value

	switch operator {
	case "and":
		return utilities.NativeBoolToBooleanObject(leftValue && rightValue)
	case "or":
		return utilities.NativeBoolToBooleanObject(leftValue || rightValue)
	case "==":
		return utilities.NativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return utilities.NativeBoolToBooleanObject(leftValue != rightValue)
	default:
		return utilities.NewError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
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
		return utilities.NativeBoolToBooleanObject(leftValue.LessThan(rightValue))
	case ">":
		return utilities.NativeBoolToBooleanObject(leftValue.GreaterThan(rightValue))
	case "<=":
		return utilities.NativeBoolToBooleanObject(leftValue.LessThanOrEqual(rightValue))
	case ">=":
		return utilities.NativeBoolToBooleanObject(leftValue.GreaterThanOrEqual(rightValue))
	case "==":
		return utilities.NativeBoolToBooleanObject(leftValue.Equal(rightValue))
	case "!=":
		return utilities.NativeBoolToBooleanObject(!leftValue.Equal(rightValue))
	case "+=":
		_, ok := env.Get(node.Left.String())

		if !ok {
			return utilities.NewError("Variable %s is unknown", node.Left.String())
		}

		dec := &object.Number{Value: leftValue.Add(rightValue)}
		env.Set(node.Left.String(), dec)

		return value.NULL
	case "-=":
		_, ok := env.Get(node.Left.String())

		if !ok {
			return utilities.NewError("Variable %s is unknown", node.Left.String())
		}

		dec := &object.Number{Value: leftValue.Sub(rightValue)}
		env.Set(node.Left.String(), dec)

		return value.NULL
	case "*=":
		_, ok := env.Get(node.Left.String())

		if !ok {
			return utilities.NewError("Variable %s is unknown", node.Left.String())
		}

		dec := &object.Number{Value: leftValue.Mul(rightValue)}
		env.Set(node.Left.String(), dec)

		return value.NULL
	case "/=":
		_, ok := env.Get(node.Left.String())

		if !ok {
			return utilities.NewError("Variable %s is unknown", node.Left.String())
		}

		dec := &object.Number{Value: leftValue.Div(rightValue)}
		env.Set(node.Left.String(), dec)

		return value.NULL
	case "..":
		numbers := make([]object.Object, 0)
		one := decimal.NewFromInt(1)
		num := leftValue

		for {
			numbers = append(numbers, &object.Number{Value: num})

			if num.GreaterThanOrEqual(rightValue) {
				break
			}

			num = num.Add(one)
		}

		return &object.List{Elements: numbers}
	default:
		return utilities.NewError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftValue + rightValue}
	case "==":
		return utilities.NativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return utilities.NativeBoolToBooleanObject(leftValue != rightValue)
	default:
		return utilities.NewError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)

	if utilities.IsError(condition) {
		return condition
	}

	if utilities.IsTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return value.NULL
	}
}

func evalIdentifierLiteral(node *ast.IdentifierLiteral, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins.Builtins[node.Value]; ok {
		return builtin
	}

	return utilities.NewError("identifier not found: " + node.Value)
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
		return utilities.NewError("index operator not supported: %s", left.Type())
	}
}

func evalListIndexExpression(list object.Object, index object.Object) object.Object {
	listObject := list.(*object.List)
	idx := index.(*object.Number).Value.IntPart()
	max := int64(len(listObject.Elements) - 1)

	if idx < 0 || idx > max {
		return value.NULL
	}

	return listObject.Elements[idx]
}

func evalMapLiteral(node *ast.MapLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.MapKey]object.MapPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)

		if utilities.IsError(key) {
			return key
		}

		mapKey, ok := key.(object.Mappable)

		if !ok {
			return utilities.NewError("unusable as map key: %s", key.Type())
		}

		val := Eval(valueNode, env)

		if utilities.IsError(val) {
			return val
		}

		mapped := mapKey.MapKey()
		pairs[mapped] = object.MapPair{Key: key, Value: val}
	}

	return &object.Map{Pairs: pairs}
}

func evalMapIndexExpression(m object.Object, index object.Object) object.Object {
	mapObject := m.(*object.Map)

	key, ok := index.(object.Mappable)

	if !ok {
		return utilities.NewError("unusable as map key: %s", index.Type())
	}

	pair, ok := mapObject.Pairs[key.MapKey()]

	if !ok {
		return value.NULL
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

		if utilities.IsError(condition) {
			return condition
		}

		if utilities.IsTruthy(condition) {
			Eval(we.Consequence, env)
		} else {
			break
		}
	}

	return value.NULL
}

func evalForExpression(fe *ast.ForExpression, env *object.Environment) object.Object {
	existingIdentifier, identifierExisted := env.Get(fe.Identifier)

	defer func() {
		if identifierExisted {
			env.Set(fe.Identifier, existingIdentifier)
		} else {
			env.Delete(fe.Identifier)
		}
	}()

	initializer := Eval(fe.Initializer, env)

	if utilities.IsError(initializer) {
		return initializer
	}

	loop := true

	for loop {
		condition := Eval(fe.Condition, env)

		if utilities.IsError(condition) {
			return condition
		}

		if utilities.IsTruthy(condition) {
			err := Eval(fe.Block, env)

			if utilities.IsError(err) {
				return err
			}

			err = Eval(fe.Increment, env)

			if utilities.IsError(err) {
				return err
			}

			continue
		}

		loop = false
	}

	return value.NULL
}

func evalForInExpression(fie *ast.ForInExpression, env *object.Environment) object.Object {
	iterable := Eval(fie.Iterable, env)

	existingKey, keyExisted := env.Get(fie.Key)
	existingValue, valueExisted := env.Get(fie.Value)

	defer func() {
		if keyExisted {
			env.Set(fie.Key, existingKey)
		} else {
			env.Delete(fie.Key)
		}

		if valueExisted {
			env.Set(fie.Value, existingValue)
		} else {
			env.Delete(fie.Value)
		}
	}()

	switch i := iterable.(type) {
	case *object.List:
		for k, v := range i.Elements {
			env.Set(fie.Key, &object.Number{Value: decimal.NewFromInt(int64(k))})
			env.Set(fie.Value, v)
			block := Eval(fie.Block, env)

			if utilities.IsError(block) {
				return block
			}
		}

		return value.NULL
	default:
		return utilities.NewError("'%s' is not a List, cannot be used in for loop", i.Inspect())
	}
}

func evalImportExpression(ie *ast.ImportExpression, env *object.Environment) object.Object {
	name := Eval(ie.Name, env)

	if utilities.IsError(name) {
		return name
	}

	if s, ok := name.(*object.String); ok {
		attributes := EvalModule(s.Value)

		if utilities.IsError(attributes) {
			return attributes
		}

		return &object.Module{Name: s.Value, Attributes: attributes}
	}

	return utilities.NewError("Import Error: invalid import path '%s'", name)
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

		return value.NULL
	default:
		return utilities.NewError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, arguments []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for key, val := range fn.Defaults {
		env.Set(key, Eval(val, env))
	}

	for index, parameter := range fn.Parameters {
		if index < len(arguments) {
			env.Set(parameter.Value, arguments[index])
		}
	}

	return env
}

func extendForEnv(fe *ast.ForExpression, forEnv *object.Environment) *object.Environment {
	env := object.NewEnclosedEnvironment(forEnv)

	env.Set(fe.Identifier, Eval(fe.Initializer, env))

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func evalAssignmentStatement(as *ast.AssignmentStatement, env *object.Environment) object.Object {
	val := Eval(as.Value, env)

	if utilities.IsError(val) {
		return val
	}

	env.Set(as.Name.Value, val)
	return nil
}
