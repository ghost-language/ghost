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

	// *Class (a Ghost-source class) is handled directly, since building an
	// instance and running its constructor is tree-walking work that only the
	// evaluator can do. Anything else claiming to be constructible —
	// currently only *object.NativeClass, a class registered by a Go program
	// embedding Ghost (§8.9, §10.3) — goes through object.Constructible
	// instead, which `new` doesn't need to know any more about than that it
	// can be asked to build an instance.
	switch class := callee.(type) {
	case *object.Class:
		arguments, err := evaluateNewArguments(node, scope)

		if err != nil {
			return err
		}

		return instantiate(class, arguments, node.Token, scope)
	case object.Constructible:
		arguments, err := evaluateNewArguments(node, scope)

		if err != nil {
			return err
		}

		return class.New(scope, node.Token, arguments...)
	default:
		return object.NewError(fault.Type, node.Token, "cannot instantiate %s, which is not a class", object.TypeName(callee)).
			WithHelp("`new` needs a class, as in `new Point(1, 2)`")
	}
}

// evaluateNewArguments evaluates a `new` expression's argument list, the one
// piece the Ghost-class and native-class construction paths need identically.
func evaluateNewArguments(node *ast.New, scope *object.Scope) ([]object.Object, object.Object) {
	arguments := evaluateExpressions(node.Arguments, scope)

	if len(arguments) == 1 && isError(arguments[0]) {
		return nil, arguments[0]
	}

	return arguments, nil
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

	result := invokeMethod(method, instance, declaringClass, arguments, caller, tok, class.Name.Value+"()")

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
	// A field initializer resolves exactly as a method body does: through the
	// scope the class was declared in, not the class's member table (§14
	// decision 12). Without this the same shadowing §13.22 describes would
	// survive here - `size = math.floor(2.7)` beside a method named `math`
	// would reach the method.
	scope := &object.Scope{Environment: class.Scope.Environment, Self: instance, Class: class, Depth: depthOf(caller)}

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

// invokeMethod runs a resolved method with the instance bound as `this` and
// the declaring class recorded so `super` resolves from the declaration
// site. name identifies the call for an arity error and for the stack frame
// recorded around a failure from running the body - never around a call that
// failed before the body ever ran (too deep, wrong arity), since a frame
// there would just repeat the position the error already reports.
func invokeMethod(method *object.Function, instance *object.Instance, declaringClass *object.Class, arguments []object.Object, caller *object.Scope, at token.Token, name string) object.Object {
	depth := depthOf(caller)

	if depth >= maxCallDepth {
		return tooDeep(at)
	}

	environment, err := createFunctionEnvironment(method, arguments, name, at)

	if err != nil {
		return err
	}

	scope := &object.Scope{Self: instance, Class: declaringClass, Environment: environment, Depth: depth + 1}

	result := Evaluate(method.Body, scope)

	if failed, ok := result.(*object.Error); ok {
		return failed.WithFrame(name, at)
	}

	return result
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
