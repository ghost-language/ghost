package modules

import (
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
	"ghostlang.org/x/ghost/value"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var KeyboardMethods = map[string]*object.LibraryFunction{}
var KeyboardProperties = map[string]*object.LibraryProperty{}

func init() {
	RegisterMethod(KeyboardMethods, "wasPressed", keyboardWasPressed)
	RegisterMethod(KeyboardMethods, "isDown", keyboardIsDown)
}

// windowCreate creates a window and OpenGL context.
func keyboardWasPressed(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 1 {
		return nil
	}

	key := int32(args[0].(*object.Number).Value.IntPart())

	if rl.IsKeyPressed(key) {
		return value.TRUE
	}

	return value.FALSE
}

// windowCreate creates a window and OpenGL context.
func keyboardIsDown(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 1 {
		return nil
	}

	key := int32(args[0].(*object.Number).Value.IntPart())

	if rl.IsKeyDown(key) {
		return value.TRUE
	}

	return value.FALSE
}
