package functions

import (
	"bufio"
	"fmt"
	"os"

	"ghostlang.org/x/ghost/object"
)

func ReadLine(args ...object.Object) object.Object {
	reader := bufio.NewReader(os.Stdin)

	if len(args) == 1 {
		prompt := args[0].(*object.String).Value + " "
		fmt.Print(prompt)
	}

	value, _ := reader.ReadString('\n')

	return &object.String{Value: string(value)}
}
