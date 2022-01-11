package modules

import (
	"path"
	"plugin"
	"strings"

	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/parser"
	"ghostlang.org/x/ghost/scanner"
	"ghostlang.org/x/ghost/token"
	"ghostlang.org/x/ghost/version"
)

var GhostMethods = map[string]*object.LibraryFunction{}
var GhostProperties = map[string]*object.LibraryProperty{}

func init() {
	RegisterMethod(GhostMethods, "abort", ghostAbort)
	RegisterMethod(GhostMethods, "execute", ghostExecute)
	RegisterMethod(GhostMethods, "extend", ghostExtend)
	RegisterMethod(GhostMethods, "identifiers", ghostIdentifiers)

	RegisterProperty(GhostProperties, "version", ghostVersion)
}

func ghostAbort(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
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

func ghostExecute(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
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

	return evaluate(program, scope)
}

func ghostExtend(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError("%d:%d: runtime error: ghost.extend() expects 1 argument. got=%d", tok.Line, tok.Column, len(args))
	}

	basePath, ok := args[0].(*object.String)

	if !ok {
		return object.NewError("%d:%d: runtime error: ghost.extend() expects the first argument to be of type 'string'. got=%s", tok.Line, tok.Column, strings.ToLower(string(args[0].Type())))
	}

	path := path.Clean(scope.Environment.GetDirectory() + "/" + basePath.Value)

	extension, err := plugin.Open(path)

	if err != nil {
		return object.NewError("%d:%d: runtime error: ghost.extend() failed opening plugin: %s", tok.Line, tok.Column, err)
	}

	register, err := extension.Lookup("Register")

	if err != nil {
		return object.NewError("%d:%d: runtime error: plugin '%s' does not contain Register function: %s", tok.Line, tok.Column, path, err)
	}

	register.(func())()

	return nil
}

func ghostIdentifiers(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewError("%d:%d: runtime error: ghost.identifiers() expects 0 arguments. got=%d", tok.Line, tok.Column, len(args))
	}

	identifiers := []object.Object{}

	store := scope.Environment.All()

	for identifier := range store {
		identifiers = append(identifiers, &object.String{Value: identifier})
	}

	return &object.List{Elements: identifiers}
}

// Properties

func ghostVersion(scope *object.Scope, tok token.Token) object.Object {
	return &object.String{Value: version.Version}
}
