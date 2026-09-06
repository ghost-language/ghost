package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

func evaluateClass(node *ast.Class, scope *object.Scope) object.Object {
	// The class keeps this scope, and its methods read through it, so it
	// outlives the block it was declared in.
	scope.Environment.Capture()

	class := &object.Class{
		Name:  node.Name,
		Scope: scope,
		Super: nil,
	}

	// super
	if node.Super != nil {
		identifier, ok := scope.Environment.Get(node.Super.Value)

		if !ok {
			return object.NewError(fault.Name, node.Super.Token, "`%s` is not defined", node.Super.Value).
				WithHelp("a class has to be declared before the class that extends it")
		}

		super, ok := identifier.(*object.Class)

		if !ok {
			return object.NewError(fault.Type, node.Super.Token, "cannot extend `%s`, which is a %s, not a class", node.Super.Value, object.TypeName(identifier))
		}

		class.Super = super
	}

	// The class environment holds the class's members and doubles as the scope
	// its body is evaluated in, so a method can call a sibling method by bare
	// name. Enclosing it in the defining scope keeps outer bindings reachable
	// from method bodies.
	class.Environment = object.NewEnclosedEnvironment(scope.Environment)
	classScope := &object.Scope{Environment: class.Environment, Self: class, Class: class}

	result := Evaluate(node.Body, classScope)

	if isError(result) {
		return result
	}

	// The same collision §13.18 rejects inside one body, reaching across two.
	// It can only be checked once the body is complete, because `use` and the
	// superclass contribute members of their own and a `use` may sit anywhere
	// in the body.
	if failed := checkInheritedCollisions(class, node); failed != nil {
		return failed
	}

	scope.Environment.Set(node.Name.Value, class)

	return class
}

// checkInheritedCollisions reports a field and a method of one name arriving
// from two different bodies - a field here against a method inherited or
// supplied by a used trait, and the reverse. Within a single body the check
// lives on the two declaration paths; neither can see this case, because at
// the moment either member is declared the other body may not have been
// consulted yet.
//
// Only field-against-method counts. A name appearing twice as two fields, or
// twice as two methods, is overriding, which is the point of extending a class
// or using a trait.
func checkInheritedCollisions(class *object.Class, node *ast.Class) object.Object {
	for _, field := range class.Fields {
		if member, owner, ok := inheritedMember(class, field.Name); ok {
			if _, isMethod := member.(*object.Function); isMethod {
				return inheritedCollisionError(field.Token, field.Name, "field", "method", owner)
			}
		}
	}

	for name, tok := range declaredMethods(node) {
		if owner, ok := inheritedField(class, name); ok {
			return inheritedCollisionError(tok, name, "method", "field", owner)
		}
	}

	return nil
}

// inheritedMember finds a member this class did not declare itself: one a
// trait supplied, or one it extends. A trait is consulted first because a
// trait used here is nearer than a superclass.
func inheritedMember(class *object.Class, name string) (object.Object, string, bool) {
	for _, trait := range class.Traits {
		if member, ok := trait.Environment.GetLocal(name); ok {
			return member, trait.Name.Value, true
		}
	}

	if class.Super != nil {
		if member, owner, ok := object.LookupMember(class.Super, name); ok {
			return member, owner.Name.Value, true
		}
	}

	return nil, "", false
}

// inheritedField answers the same question for fields, which live on their
// declarer rather than in an environment and so need their own walk.
func inheritedField(class *object.Class, name string) (string, bool) {
	for _, trait := range class.Traits {
		if trait.HasField(name) {
			return trait.Name.Value, true
		}
	}

	for ancestor := class.Super; ancestor != nil; ancestor = ancestor.Super {
		for _, trait := range ancestor.Traits {
			if trait.HasField(name) {
				return trait.Name.Value, true
			}
		}

		if ancestor.HasField(name) {
			return ancestor.Name.Value, true
		}
	}

	return "", false
}

// declaredMethods reads the method names this class body declares, with the
// token each was written at. The evaluated class holds the functions but not
// where they came from, and a collision deserves to be reported at the
// declaration that created it.
func declaredMethods(node *ast.Class) map[string]token.Token {
	methods := map[string]token.Token{}

	for _, statement := range node.Body.Statements {
		if function, ok := methodOf(statement); ok {
			methods[function.Name.Value] = function.Name.Token
		}
	}

	return methods
}

func methodOf(statement ast.StatementNode) (*ast.Function, bool) {
	switch node := statement.(type) {
	case *ast.Function:
		if node.Name != nil {
			return node, true
		}
	case *ast.Expression:
		if function, ok := node.Expression.(*ast.Function); ok && function.Name != nil {
			return function, true
		}
	}

	return nil, false
}
