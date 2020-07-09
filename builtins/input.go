package builtins

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"ghostlang.org/ghost/object"
)

func Input(args ...object.Object) object.Object {
	if len(args) == 1 {
		prompt := args[0].(*object.String).Value + " "
		fmt.Fprintf(os.Stdout, prompt)
	}

	buffer := bufio.NewReader(os.Stdin)

	line, _, err := buffer.ReadLine()

	if err != nil && err != io.EOF {
		return newError(fmt.Sprintf("error reading input: %s", err))
	}
	return &object.String{Value: string(line)}
}
