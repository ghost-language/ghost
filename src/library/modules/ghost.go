package modules

import (
	"strings"

	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/parser"
	"ghostlang.org/x/ghost/scanner"
	"ghostlang.org/x/ghost/token"
	"ghostlang.org/x/ghost/version"
)

var Ghost = map[string]*object.LibraryFunction{}

func init() {
	RegisterMethod(Ghost, "abort", ghostAbort)
	RegisterMethod(Ghost, "execute", ghostExecute)
	RegisterMethod(Ghost, "identifiers", ghostIdentifiers)
	RegisterMethod(Ghost, "version", ghostVersion)
}

func ghostAbort(env *object.Environment, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError("%d:%d: runtime error: ghost.abort() expects 1 argument. got=%d", tok.Line, tok.Column, len(args))
	}

	switch obj := args[0].(type) {
	case *object.Null:
		return nil
	case *object.String:
		return object.NewError(obj.Value)
	}

	return object.NewError("%d:%d: runtime error: ghost.abort() expects the first argument to be of type 'null' or 'string'. got=%s", tok.Line, tok.Column, strings.ToLower(string(args[0].Type())))
}

func ghostExecute(env *object.Environment, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError("%d:%d: runtime error: ghost.execute() expects 1 argument. got=%d", tok.Line, tok.Column, len(args))
	}

	source, ok := args[0].(*object.String)

	if !ok {
		return object.NewError("%d:%d: runtime error: ghost.execute() expects the first argument to be of type 'string'. got=%s", tok.Line, tok.Column, strings.ToLower(string(args[0].Type())))
	}

	scanner := scanner.New(source.Value)
	parser := parser.New(scanner)
	program := parser.Parse()

	return evaluate(program, env)
}

func ghostIdentifiers(env *object.Environment, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewError("%d:%d: runtime error: ghost.identifiers() expects 0 arguments. got=%d", tok.Line, tok.Column, len(args))
	}

	identifiers := []object.Object{}

	store := env.All()

	for identifier := range store {
		identifiers = append(identifiers, &object.String{Value: identifier})
	}

	return &object.List{Elements: identifiers}
}

func ghostVersion(env *object.Environment, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewError("%d:%d: runtime error: ghost.version() expects 0 arguments. got=%d", tok.Line, tok.Column, len(args))
	}

	return &object.String{Value: version.Version}
}
