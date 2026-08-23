package object

import "ghostlang.org/x/ghost/token"

// Constructible is implemented by a value `new` can build an instance from
// that is not a Ghost-source *Class — in practice, *NativeClass. The
// evaluator handles *Class directly (constructing an *Instance and running
// its constructor is tree-walking work that belongs there, not here); this
// interface is `new`'s other path, for a class whose instances are built and
// driven entirely by Go code.
type Constructible interface {
	Object
	New(scope *Scope, tok token.Token, args ...Object) Object
}
