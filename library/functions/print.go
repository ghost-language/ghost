package functions

import (
	"fmt"
	"strings"

	"ghostlang.org/x/ghost/log"
	"ghostlang.org/x/ghost/object"
)

func Print(args ...object.Object) object.Object {
	if len(args) > 0 {
		str := make([]string, 0)

		for _, value := range args {
			str = append(str, value.String())
		}

		log.Info(strings.Join(str, " "))
	} else {
		fmt.Println()
	}

	return nil
}
