package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/fault"
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
		return object.NewError(fault.Type, node.Token, "cannot instantiate a null value").
			WithHelp("`new` needs a class, as in `new Point(1, 2)`")
	}

	class, ok := callee.(*object.Class)

	if !ok {
		return object.NewError(fault.Type, node.Token, "cannot instantiate %s, which is not a class", object.TypeName(callee)).
			WithHelp("`new` needs a class, as in `new Point(1, 2)`")
	}

	arguments := evaluateExpressions(node.Arguments, scope)

	if len(arguments) == 1 && isError(arguments[0]) {
		return arguments[0]
	}

	return instantiate(class, arguments, node.Token, scope)
}

// instantiate builds an instance of the class, initializes its fields, and runs
// the nearest constructor in the class chain.
func instantiate(class *object.Class, arguments []object.Object, tok token.Token, caller *object.Scope) object.Object {
	instance := &object.Instance{Class: class, Environment: object.NewEnclosedEnvironment(class.Environment)}

	if result := initializeFields(instance, class, caller); result != nil {
		return result
	}

	constructor, declaringClass, ok := instance.LookupMember(constructorName)

	if !ok {
		return instance
	}

	method, ok := constructor.(*object.Function)

	if !ok {
		return object.NewError(fault.Type, tok, "the constructor of class `%s` is a %s, not a function", class.Name.Value, object.TypeName(constructor))
	}

	result := invokeMethod(method, instance, declaringClass, arguments, caller, tok)

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
func initializeFields(instance *object.Instance, class *object.Class, caller *object.Scope) object.Object {
	for _, ancestor := range class.Ancestors() {
		for _, trait := range ancestor.Traits {
			for _, field := range trait.Fields {
				if result := initializeField(instance, ancestor, field, caller); result != nil {
					return result
				}
			}
		}

		for _, field := range ancestor.Fields {
			if result := initializeField(instance, ancestor, field, caller); result != nil {
				return result
			}
		}
	}

	return nil
}

func initializeField(instance *object.Instance, class *object.Class, field object.Field, caller *object.Scope) object.Object {
	scope := &object.Scope{Environment: class.Environment, Self: instance, Class: class, Depth: depthOf(caller)}

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
func invokeMethod(method *object.Function, instance *object.Instance, declaringClass *object.Class, arguments []object.Object, caller *object.Scope, at token.Token) object.Object {
	depth := depthOf(caller)

	if depth >= maxCallDepth {
		return tooDeep(at)
	}

	environment := createFunctionEnvironment(method, arguments)
	scope := &object.Scope{Self: instance, Class: declaringClass, Environment: environment, Depth: depth + 1}

	return Evaluate(method.Body, scope)
}

// depthOf reads how deep a caller already is. A method reached from Go rather
// than from Ghost — an embedder calling into a script — has no caller, and
// starts at the bottom.
func depthOf(caller *object.Scope) int {
	if caller == nil {
		return 0
	}

	return caller.Depth
}
