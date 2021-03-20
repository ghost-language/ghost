package standard

import (
	"ghostlang.org/x/ghost/object"
)

var StandardFunctions = map[string]*object.Standard{}

func RegisterFunction(name string, function object.StandardFunction) {
	StandardFunctions[name] = &object.Standard{Function: function}
}
