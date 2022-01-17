package modules

import (
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var DrawMethods = map[string]*object.LibraryFunction{}
var DrawProperties = map[string]*object.LibraryProperty{}

func init() {
	RegisterMethod(DrawMethods, "begin", drawBegin)
	RegisterMethod(DrawMethods, "end", drawEnd)
	RegisterMethod(DrawMethods, "clearBackground", drawClearBackground)
	RegisterMethod(DrawMethods, "rectangle", drawRectangle)
}

// windowCreate creates a window and OpenGL context.
func drawBegin(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	rl.BeginDrawing()

	return nil
}

// windowCreate creates a window and OpenGL context.
func drawEnd(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	rl.EndDrawing()

	return nil
}

// windowCreate creates a window and OpenGL context.
func drawClearBackground(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	rl.ClearBackground(rl.RayWhite)

	return nil
}

func drawRectangle(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 4 {
		return nil
	}

	xPosition := int32(args[0].(*object.Number).Value.IntPart())
	yPosition := int32(args[1].(*object.Number).Value.IntPart())
	width := int32(args[2].(*object.Number).Value.IntPart())
	height := int32(args[3].(*object.Number).Value.IntPart())

	rl.DrawRectangle(xPosition, yPosition, width, height, rl.Black)

	return nil
}
