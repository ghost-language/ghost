package modules

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

var OsMethods = map[string]*object.LibraryFunction{}
var OsProperties = map[string]*object.LibraryProperty{}

func init() {
	RegisterMethod(OsMethods, "args", osArgs)
	RegisterMethod(OsMethods, "exit", osExit)
	RegisterMethod(OsMethods, "sleep", osSleep)

	RegisterProperty(OsProperties, "name", osName)
}

func osArgs(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("os.args", tok, args, 0); err != nil {
		return err
	}

	arguments := os.Args[1:]
	elements := make([]object.Object, len(arguments))

	for index, argument := range arguments {
		elements[index] = &object.String{Value: argument}
	}

	return &object.List{Elements: elements}
}

// osExit ends the program with a status code, and optionally a parting message.
// It is the one library method that never returns, so its arguments are checked
// before anything is printed: exiting on a miscall would take the mistake with
// it.
func osExit(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arityRange("os.exit", tok, args, 1, 2); err != nil {
		return err
	}

	status, err := integerAt("os.exit", tok, args, 0)

	if err != nil {
		return err
	}

	if len(args) == 2 {
		message, err := stringAt("os.exit", tok, args, 1)

		if err != nil {
			return err
		}

		fmt.Println(message)
	}

	os.Exit(int(status))

	return nil
}

// osSleep pauses the running program for a duration in milliseconds. It is a
// process control, the same family as exit - pausing or ending the program
// itself - rather than a date operation, which is why it lives here and not
// in the date module.
func osSleep(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("os.sleep", tok, args, 1); err != nil {
		return err
	}

	milliseconds, err := integerAt("os.sleep", tok, args, 0)

	if err != nil {
		return err
	}

	if milliseconds < 0 {
		return object.NewError(fault.Value, tok, "`os.sleep()` expects a duration of zero or greater, got %d", milliseconds)
	}

	time.Sleep(time.Duration(milliseconds) * time.Millisecond)

	return nil
}

// Properties

func osName(scope *object.Scope, tok token.Token) object.Object {
	return &object.String{Value: runtime.GOOS}
}
