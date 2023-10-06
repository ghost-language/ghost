package modules

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
	"ghostlang.org/x/ghost/value"
)

var ConsoleMethods = map[string]*object.LibraryFunction{}
var ConsoleProperties = map[string]*object.LibraryProperty{}

func init() {
	RegisterMethod(ConsoleMethods, "error", consoleError)
	RegisterMethod(ConsoleMethods, "info", consoleInfo)
	RegisterMethod(ConsoleMethods, "log", consoleLog)
	RegisterMethod(ConsoleMethods, "read", consoleRead)
	RegisterMethod(ConsoleMethods, "warn", consoleWarn)
	RegisterMethod(ConsoleMethods, "clear", consoleClear)
	RegisterMethod(ConsoleMethods, "print", consolePrint)
	RegisterMethod(ConsoleMethods, "newLine", consoleNewLine)
}

func consoleError(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	values := make([]string, 0)

	for _, value := range args {
		values = append(values, value.String())
	}

	printLine(values, "error")

	return nil
}

func consoleInfo(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	values := make([]string, 0)

	for _, value := range args {
		values = append(values, value.String())
	}

	printLine(values, "info")

	return nil
}

func consoleLog(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	values := make([]string, 0)

	for _, value := range args {
		values = append(values, value.String())
	}

	printLine(values, "")

	return nil
}

func consoleRead(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	scanner := bufio.NewScanner(os.Stdin)

	if len(args) == 1 {
		prompt := args[0].(*object.String).Value

		fmt.Print(prompt)
	}

	val := scanner.Scan()

	if !val {
		return value.NULL
	}

	return &object.String{Value: scanner.Text()}
}

func consoleWarn(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	values := make([]string, 0)

	for _, value := range args {
		values = append(values, value.String())
	}

	printLine(values, "warning")

	return nil
}

func consoleClear(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	return nil
}

func consolePrint(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	values := make([]string, 0)

	for _, value := range args {
		values = append(values, value.String())
	}

	print(values)

	return nil
}

func consoleNewLine(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	fmt.Println()

	return nil
}

//

func printLine(values []string, prefix string) {
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

func print(values []string) {
	if len(values) > 0 {
		str := make([]string, 0)

		str = append(str, values...)

		fmt.Print(strings.Join(str, " "))
	}
}
