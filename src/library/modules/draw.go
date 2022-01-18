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
	RegisterMethod(DrawMethods, "rectangleOutline", drawRectangleOutline)
	RegisterMethod(DrawMethods, "circle", drawCircle)
	RegisterMethod(DrawMethods, "circleOutline", drawCircleOutline)
	RegisterMethod(DrawMethods, "line", drawLine)
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

func drawRectangleOutline(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 4 {
		return nil
	}

	xPosition := int32(args[0].(*object.Number).Value.IntPart())
	yPosition := int32(args[1].(*object.Number).Value.IntPart())
	width := int32(args[2].(*object.Number).Value.IntPart())
	height := int32(args[3].(*object.Number).Value.IntPart())

	rl.DrawRectangleLines(xPosition, yPosition, width, height, rl.Black)

	return nil
}

func drawCircle(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 3 {
		return nil
	}

	xCenter := int32(args[0].(*object.Number).Value.IntPart())
	yCenter := int32(args[1].(*object.Number).Value.IntPart())
	radius := float32(args[2].(*object.Number).Value.InexactFloat64())

	rl.DrawCircle(xCenter, yCenter, radius, rl.Black)

	return nil
}

func drawCircleOutline(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 3 {
		return nil
	}

	xCenter := int32(args[0].(*object.Number).Value.IntPart())
	yCenter := int32(args[1].(*object.Number).Value.IntPart())
	radius := float32(args[2].(*object.Number).Value.InexactFloat64())

	rl.DrawCircleLines(xCenter, yCenter, radius, rl.Black)

	return nil
}

func drawLine(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 4 {
		return nil
	}

	xStart := int32(args[0].(*object.Number).Value.IntPart())
	yStart := int32(args[1].(*object.Number).Value.IntPart())
	xEnd := int32(args[2].(*object.Number).Value.IntPart())
	yEnd := int32(args[3].(*object.Number).Value.IntPart())

	rl.DrawLine(xStart, yStart, xEnd, yEnd, rl.Black)

	return nil
}
