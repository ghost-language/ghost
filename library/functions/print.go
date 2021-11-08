package functions

import (
	"fmt"
	"strings"

	"ghostlang.org/x/ghost/object"
)

func Print(args ...object.Object) object.Object {
	if len(args) > 0 {
		str := make([]string, 0)

		for _, value := range args {
			str = append(str, value.String())
		}

		fmt.Println(strings.Join(str, " "))
	} else {
		fmt.Println()
	}

	return nil
}
