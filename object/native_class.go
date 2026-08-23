package object

import "ghostlang.org/x/ghost/token"

// NativeClass is a class whose instances are built and driven entirely by Go
// code, rather than by evaluating a Ghost class body — the counterpart to
// *Class for a Go program embedding Ghost that wants to hand scripts a
// stateful host resource (an audio handle, a window, ...) as something they
// `new` and call methods on, the way `new ClassName()` already works for any
// Ghost-defined class, without paying to tree-walk methods that only ever
// call back into Go anyway.
//
// A host registers one under a scheme with library.RegisterClassForScheme
// (or the unscoped library.RegisterClass, for the standard library's own
// "ghost" scheme), which is how `import { Audio } from "lumen:audio"` and
// `new Audio(path)` both resolve — Constructor builds and returns whatever
// object Audio's own instances should be; that returned value owns deciding
// what its Method() does, entirely on the host's side.
type NativeClass struct {
	Name        string
	Constructor GoFunction
}

// String represents the class object's value as a string. Deliberately the
// same as *Class's — a script has no reason to know or care whether a class
// it is using is implemented in Ghost or in Go.
func (class *NativeClass) String() string {
	return "class"
}

// Type returns CLASS, the same type *Class reports, for the same reason:
// `type(Audio)` should read "class" regardless of which kind of class it is.
func (class *NativeClass) Type() Type {
	return CLASS
}

// Method defines the set of methods available on the class object itself, as
// opposed to its instances — there are none, matching *Class: a class value
// only supports `new`.
func (class *NativeClass) Method(method string, tok token.Token, args []Object) (Object, bool) {
	return nil, false
}

// New builds a new instance by calling the registered constructor, the same
// way calling any other library function does.
func (class *NativeClass) New(scope *Scope, tok token.Token, args ...Object) Object {
	return class.Constructor(scope, tok, args...)
}
