package object

import "ghostlang.org/x/ghost/token"

// Scope objects consist of an environment and parent object.
type Scope struct {
	Environment *Environment
	Self        Object

	// Depth is how many calls deep this scope is. It rides on the scope rather
	// than on a counter in the evaluator because Ghost code can run on more
	// than one goroutine at a time — an http.handle() callback runs per request
	// — and a shared counter would both race and mean the wrong thing, charging
	// one request's recursion against another's budget.
	Depth int

	// Class records which class declared the code currently executing. It is
	// what `super` resolves against, so an inherited method always starts its
	// super lookup from its own declaring class rather than from the receiver's
	// class — otherwise a method inherited two levels down would resolve
	// `super` back to itself.
	Class *Class
}

// Enclose returns a scope for a block to run in: a fresh environment chained
// to this one, and everything else about the scope — the receiver, the
// declaring class, the call depth — carried over unchanged.
//
// The scope released here last is reused when nothing captured it, so a block
// inside a hot loop costs no allocation per execution. What makes that safe is
// Environment.Capture: a closure, class, or trait created inside the block
// marks the environment, and a marked one is dropped rather than reused.
func (scope *Scope) Enclose() *Scope {
	child := scope.Environment.freeChild

	if child == nil {
		return &Scope{
			Environment: NewEnclosedEnvironment(scope.Environment),
			Self:        scope.Self,
			Depth:       scope.Depth,
			Class:       scope.Class,
		}
	}

	scope.Environment.freeChild = nil

	child.Environment.clear()
	child.Self = scope.Self
	child.Depth = scope.Depth
	child.Class = scope.Class

	return child
}

// Release offers a finished block scope back for the next block to reuse. A
// scope whose environment something captured is dropped instead, so the value
// that captured it keeps reading the bindings it closed over.
//
// Not calling it is only ever a missed reuse, never a correctness problem.
func (scope *Scope) Release(child *Scope) {
	if child.Environment.captured {
		return
	}

	scope.Environment.freeChild = child
}

// String represents the scope object's value as a string.
func (scope *Scope) String() string {
	return "scope"
}

// Type returns the scope object type.
func (scope *Scope) Type() Type {
	return SCOPE
}

// Method defines the set of methods available on scope objects.
func (scope *Scope) Method(method string, tok token.Token, args []Object) (Object, bool) {
	return nil, false
}
