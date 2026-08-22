package modules

import (
	"path"
	"plugin"

	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/optimizer"
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

// ghostAbort stops the program with a message of the script's own choosing. It
// is the one place where a Ghost program raises an error rather than receiving
// one, so the message is taken verbatim and reported as a value error: the
// program decided its own state was wrong.
func ghostAbort(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("ghost.abort", tok, args, 1); err != nil {
		return err
	}

	switch reason := args[0].(type) {
	case *object.Null:
		return nil
	case *object.String:
		return object.NewError(fault.Value, tok, "%s", reason.Value)
	}

	return object.NewError(fault.Argument, tok, "`ghost.abort()` expects argument 1 to be a string or null, got %s", object.TypeName(args[0]))
}

func ghostExecute(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("ghost.execute", tok, args, 1); err != nil {
		return err
	}

	code, err := stringAt("ghost.execute", tok, args, 0)

	if err != nil {
		return err
	}

	parsed := parser.New(scanner.New(code, tok.File))
	program := parsed.Parse()

	// Code handed to `ghost.execute` is parsed on its own, so its syntax errors
	// have to be reported here rather than by whoever parsed the outer script.
	if raised := parsed.Errors(); len(raised) != 0 {
		return object.NewErrorFrom(raised[0])
	}

	return evaluate(optimizer.Optimize(program), scope)
}

func ghostExtend(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("ghost.extend", tok, args, 1); err != nil {
		return err
	}

	target, err := stringAt("ghost.extend", tok, args, 0)

	if err != nil {
		return err
	}

	resolved := path.Clean(scope.Environment.GetDirectory() + "/" + target)

	extension, failure := plugin.Open(resolved)

	if failure != nil {
		return systemFailure("ghost.extend", tok, failure)
	}

	register, failure := extension.Lookup("Register")

	if failure != nil {
		return object.NewError(fault.System, tok, "plugin `%s` has no Register function", resolved).
			WithHelp("a plugin has to export `func Register()` for Ghost to call")
	}

	entry, ok := register.(func())

	if !ok {
		return object.NewError(fault.System, tok, "the Register function in plugin `%s` has the wrong signature", resolved).
			WithHelp("a plugin has to export `func Register()`, taking no arguments and returning nothing")
	}

	entry()

	return nil
}

func ghostIdentifiers(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("ghost.identifiers", tok, args, 0); err != nil {
		return err
	}

	store := scope.Environment.All()
	identifiers := make([]object.Object, 0, len(store))

	for identifier := range store {
		identifiers = append(identifiers, &object.String{Value: identifier})
	}

	return &object.List{Elements: identifiers}
}

// Properties

func ghostVersion(scope *object.Scope, tok token.Token) object.Object {
	return &object.String{Value: version.Version}
}
