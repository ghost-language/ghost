package modules

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"ghostlang.org/x/ghost/color"
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
	if err := arityRange("console.read", tok, args, 0, 1); err != nil {
		return err
	}

	scanner := bufio.NewScanner(os.Stdin)

	if len(args) == 1 {
		prompt, err := stringAt("console.read", tok, args, 0)

		if err != nil {
			return err
		}

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

// printLine writes a console line, colouring its prefix by what the line is for.
// The prefix says how to read what follows, so it is styled the same way the
// rest of Ghost styles that meaning — and dropped entirely when the output is
// not going to a terminal that can show it.
func printLine(values []string, prefix string) {
	if len(values) == 0 {
		fmt.Println()

		return
	}

	parts := make([]string, 0, len(values)+1)

	if prefix != "" {
		parts = append(parts, style(prefix)+":")
	}

	parts = append(parts, values...)

	fmt.Println(strings.Join(parts, " "))
}

// style paints a console prefix in the colour its meaning has everywhere else.
func style(prefix string) string {
	profile := color.Detect(os.Stdout)

	switch prefix {
	case "error":
		return profile.Error(prefix)
	case "warning":
		return profile.Warning(prefix)
	case "info":
		return profile.Debug(prefix)
	}

	return prefix
}

func print(values []string) {
	if len(values) > 0 {
		str := make([]string, 0)

		str = append(str, values...)

		fmt.Print(strings.Join(str, " "))
	}
}
