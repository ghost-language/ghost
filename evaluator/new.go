package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
	"ghostlang.org/x/ghost/value"
)

// constructorName is the method a class runs when instantiated.
const constructorName = "constructor"

func evaluateNew(node *ast.New, scope *object.Scope) object.Object {
	callee := Evaluate(node.Class, scope)

	if isError(callee) {
		return callee
	}

	if callee == nil {
		return newError("%d:%d:%s: runtime error: cannot instantiate a non-class value", node.Token.Line, node.Token.Column, node.Token.File)
	}

	class, ok := callee.(*object.Class)

	if !ok {
		return newError("%d:%d:%s: runtime error: cannot instantiate a non-class value, got %s", node.Token.Line, node.Token.Column, node.Token.File, callee.Type())
	}

	arguments := evaluateExpressions(node.Arguments, scope)

	if len(arguments) == 1 && isError(arguments[0]) {
		return arguments[0]
	}

	return instantiate(class, arguments, node.Token)
}

// instantiate builds an instance of the class, initializes its fields, and runs
// the nearest constructor in the class chain.
func instantiate(class *object.Class, arguments []object.Object, tok token.Token) object.Object {
	instance := &object.Instance{Class: class, Environment: object.NewEnclosedEnvironment(class.Environment)}

	if result := initializeFields(instance, class); result != nil {
		return result
	}

	constructor, declaringClass, ok := instance.LookupMember(constructorName)

	if !ok {
		return instance
	}

	method, ok := constructor.(*object.Function)

	if !ok {
		return newError("%d:%d:%s: runtime error: constructor of class %s is not a function, got %s", tok.Line, tok.Column, tok.File, class.Name.Value, constructor.Type())
	}

	result := invokeMethod(method, instance, declaringClass, arguments)

	if isError(result) {
		return result
	}

	return instance
}

// initializeFields evaluates each field declaration into the instance itself,
// so every instance owns its values rather than sharing one copy held on the
// class. Ancestors are applied first so a subclass field overrides the
// declaration it shadows, and a class's own fields are applied after the traits
// it uses for the same reason.
func initializeFields(instance *object.Instance, class *object.Class) object.Object {
	for _, ancestor := range class.Ancestors() {
		for _, trait := range ancestor.Traits {
			for _, field := range trait.Fields {
				if result := initializeField(instance, ancestor, field); result != nil {
					return result
				}
			}
		}

		for _, field := range ancestor.Fields {
			if result := initializeField(instance, ancestor, field); result != nil {
				return result
			}
		}
	}

	return nil
}

func initializeField(instance *object.Instance, class *object.Class, field object.Field) object.Object {
	scope := &object.Scope{Environment: class.Environment, Self: instance, Class: class}

	result := Evaluate(field.Value, scope)

	if isError(result) {
		return result
	}

	if result == nil {
		result = value.NULL
	}

	instance.Environment.Set(field.Name, result)

	return nil
}

// invokeMethod runs a resolved method with the instance bound as `this` and the
// declaring class recorded so `super` resolves from the declaration site.
func invokeMethod(method *object.Function, instance *object.Instance, declaringClass *object.Class, arguments []object.Object) object.Object {
	environment := createFunctionEnvironment(method, arguments)
	scope := &object.Scope{Self: instance, Class: declaringClass, Environment: environment}

	return Evaluate(method.Body, scope)
}
