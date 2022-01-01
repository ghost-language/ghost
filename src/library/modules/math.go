package modules

import (
	"math/rand"
	"time"

	"ghostlang.org/x/ghost/object"
	"github.com/shopspring/decimal"
)

var Math = map[string]*object.LibraryFunction{}

func init() {
	Math = RegisterMethod(Math, "pi", mathPi)
	Math = RegisterMethod(Math, "random", mathRandom)
	Math = RegisterMethod(Math, "seed", mathSeed)
}

func mathPi(args ...object.Object) object.Object {
	pi, _ := decimal.NewFromString("3.14159265358979323846264338327950288419716939937510582097494459")

	return &object.Number{Value: pi}
}

func mathRandom(args ...object.Object) object.Object {
	min := float64(0)
	max := float64(0)

	if len(args) > 0 {
		max, _ = args[0].(*object.Number).Value.Float64()

		if len(args) > 1 {
			min = max
			max, _ = args[1].(*object.Number).Value.Float64()
		}
	}

	number := float64(0)

	if max > 0 {
		number = float64(min + rand.Float64()*(max-min))
	} else {
		number = rand.Float64()
	}

	return &object.Number{Value: decimal.NewFromFloat(number)}
}

func mathSeed(args ...object.Object) object.Object {
	var seed int64

	if len(args) == 1 && args[0].Type() == object.NUMBER {
		seed = args[0].(*object.Number).Value.IntPart()
	} else {
		seed = time.Now().UnixNano()
	}

	rand.Seed(seed)

	return nil
}
