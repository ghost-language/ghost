package functions

import (
	"fmt"
	"strings"

	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

func Print(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if len(args) > 0 {
		str := make([]string, 0)

		for _, value := range args {
			str = append(str, value.String())
		}

		fmt.Fprintln(scope.Environment.GetWriter(), strings.Join(str, " "))
	} else {
		fmt.Fprintln(scope.Environment.GetWriter())
	}

	return nil
}
