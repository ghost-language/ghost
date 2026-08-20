package ast

import "ghostlang.org/x/ghost/token"

// Library binding states for an identifier.
//
// Resolving an identifier means consulting the library module and function
// registries before the scope chain, which costs two map lookups keyed by
// string. Ordinary variables miss both every time they are read, and in a loop
// that hashing dominated identifier evaluation.
//
// The optimizer classifies each identifier once, before evaluation begins, so
// the evaluator can skip the registries for names that cannot be library
// globals. The zero value is deliberately "unknown": an AST that was never
// optimized keeps the original behavior of always consulting the registries.
//
// This is written during optimization and only read during evaluation, which
// matters because Ghost code can be evaluated concurrently - http.handle()
// callbacks run per request on their own goroutines.
const (
	LibraryBindingUnknown uint8 = iota // not analyzed; consult the registries
	LibraryBindingLocal                // cannot be a library global
	LibraryBindingGlobal               // names a library module or function
)

type Identifier struct {
	AssignmentNode
	Token token.Token
	Value string

	// LibraryBinding is set by the optimizer; see the constants above.
	LibraryBinding uint8
}
