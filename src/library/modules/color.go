package modules

import (
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

var ColorMethods = map[string]*object.LibraryFunction{}
var ColorProperties = map[string]*object.LibraryProperty{}

func init() {
	RegisterMethod(ColorMethods, "new", colorNew)
}

func colorNew(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 4 {
		return object.NewError("%d:%d: runtime error: color.new() expects 4 arguments. got=%d", tok.Line, tok.Column, len(args))
	}

	color := make(map[string]interface{})

	color["r"] = int32(args[0].(*object.Number).Value.IntPart())
	color["g"] = int32(args[1].(*object.Number).Value.IntPart())
	color["b"] = int32(args[2].(*object.Number).Value.IntPart())
	color["a"] = int32(args[3].(*object.Number).Value.IntPart())

	return object.NewMap(color)
}
