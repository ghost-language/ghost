package functions

import (
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

// Type names the type of a value. The name it answers with is the same name
// error messages use, so a reader who has been told a value is a `list` can
// test for exactly that.
func Type(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := object.Arity("type()", tok, args, 1); err != nil {
		return err
	}

	return &object.String{Value: object.TypeName(args[0])}
}
