package modules

import (
	"fmt"
	"strings"

	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
	"github.com/peterh/liner"
)

var ConsoleMethods = map[string]*object.LibraryFunction{}
var ConsoleProperties = map[string]*object.LibraryProperty{}

func init() {
	RegisterMethod(ConsoleMethods, "error", consoleError)
	RegisterMethod(ConsoleMethods, "info", consoleInfo)
	RegisterMethod(ConsoleMethods, "log", consoleLog)
	RegisterMethod(ConsoleMethods, "read", consoleRead)
	RegisterMethod(ConsoleMethods, "warn", consoleWarn)
}

func consoleError(env *object.Environment, tok token.Token, args ...object.Object) object.Object {
	values := make([]string, 0)

	for _, value := range args {
		values = append(values, value.String())
	}

	print(values, "error")

	return nil
}

func consoleInfo(env *object.Environment, tok token.Token, args ...object.Object) object.Object {
	values := make([]string, 0)

	for _, value := range args {
		values = append(values, value.String())
	}

	print(values, "info")

	return nil
}

func consoleLog(env *object.Environment, tok token.Token, args ...object.Object) object.Object {
	values := make([]string, 0)

	for _, value := range args {
		values = append(values, value.String())
	}

	print(values, "")

	return nil
}

func consoleRead(env *object.Environment, tok token.Token, args ...object.Object) object.Object {
	line := liner.NewLiner()
	prompt := ""
	defer line.Close()

	if len(args) == 1 {
		prompt = args[0].(*object.String).Value + " "
	}

	value, _ := line.Prompt(prompt)

	return &object.String{Value: string(value)}
}

func consoleWarn(env *object.Environment, tok token.Token, args ...object.Object) object.Object {
	values := make([]string, 0)

	for _, value := range args {
		values = append(values, value.String())
	}

	print(values, "warning")

	return nil
}

//

func print(values []string, prefix string) {
	if len(values) > 0 {
		str := make([]string, 0)

		if len(prefix) > 0 {
			str = append(str, prefix+":")
		}

		str = append(str, values...)

		fmt.Println(strings.Join(str, " "))
	} else {
		fmt.Println()
	}
}
