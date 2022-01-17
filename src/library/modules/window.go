package modules

import (
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
	"ghostlang.org/x/ghost/value"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/shopspring/decimal"
)

var WindowMethods = map[string]*object.LibraryFunction{}
var WindowProperties = map[string]*object.LibraryProperty{}

func init() {
	rl.SetTraceLog(rl.LogNone)

	RegisterMethod(WindowMethods, "create", windowCreate)
	RegisterMethod(WindowMethods, "shouldClose", windowShouldClose)
	RegisterMethod(WindowMethods, "close", windowClose)
	RegisterMethod(WindowMethods, "setTargetFPS", windowSetTargetFPS)
	RegisterMethod(WindowMethods, "toggleFullscreen", windowToggleFullscreen)

	RegisterProperty(WindowProperties, "FPS", windowFPS)
	RegisterProperty(WindowProperties, "width", windowWidth)
	RegisterProperty(WindowProperties, "height", windowHeight)
}

// windowCreate creates a window and OpenGL context.
func windowCreate(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 3 {
		return nil
	}

	width := int32(args[0].(*object.Number).Value.IntPart())
	height := int32(args[1].(*object.Number).Value.IntPart())
	title := args[2].(*object.String).Value

	rl.InitWindow(width, height, title)

	return nil
}

// windowShouldClose checks if KEY_ESCAPED or close icon was pressed.
func windowShouldClose(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if rl.WindowShouldClose() {
		return value.TRUE
	}

	return value.FALSE
}

// windowClose closes the window and unloads OpenGL context.
func windowClose(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	rl.CloseWindow()

	return nil
}

func windowSetTargetFPS(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 1 {
		return nil
	}

	fps := int32(args[0].(*object.Number).Value.IntPart())

	rl.SetTargetFPS(fps)

	return nil
}

func windowToggleFullscreen(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	rl.ToggleFullscreen()

	return nil
}

func windowFPS(scope *object.Scope, tok token.Token) object.Object {
	val := decimal.NewFromFloat32(rl.GetFPS())

	return &object.Number{Value: val}
}

func windowWidth(scope *object.Scope, tok token.Token) object.Object {
	val := decimal.NewFromInt(int64(rl.GetScreenWidth()))

	return &object.Number{Value: val}
}

func windowHeight(scope *object.Scope, tok token.Token) object.Object {
	val := decimal.NewFromInt(int64(rl.GetScreenHeight()))

	return &object.Number{Value: val}
}
