package modules

import (
	"math/rand"
	"time"

	"ghostlang.org/x/ghost/log"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

var Math = []*object.LibraryFunction{}

func init() {
	log.Debug("calling math module...")

	rand.Seed(time.Now().UnixNano())

	Math = RegisterMethod(Math, "pi", mathPi)
}

func mathPi(args ...object.Object) object.Object {
	return value.NULL
}
