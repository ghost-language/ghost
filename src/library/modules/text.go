package modules

import (
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var TextMethods = map[string]*object.LibraryFunction{}
var TextProperties = map[string]*object.LibraryProperty{}

func init() {
	RegisterMethod(TextMethods, "draw", textDraw)
}

// windowCreate creates a window and OpenGL context.
func textDraw(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 4 {
		return nil
	}

	text := args[0].(*object.String).Value
	xPosition := int32(args[1].(*object.Number).Value.IntPart())
	yPosition := int32(args[2].(*object.Number).Value.IntPart())
	size := int32(args[3].(*object.Number).Value.IntPart())
	color := rl.Black

	rl.DrawText(text, xPosition, yPosition, size, color)

	return nil
}
