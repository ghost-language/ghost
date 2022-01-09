package modules

import (
	"math/rand"
	"time"

	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
	"github.com/shopspring/decimal"
)

var RandomMethods = map[string]*object.LibraryFunction{}
var RandomProperties = map[string]*object.LibraryProperty{}
var seed int64 = 0

func init() {
	// Set an initial fixed seed value. This can be overwritten at runtime.
	rand.Seed(seed)

	RegisterMethod(RandomMethods, "seed", randomSeed)
	RegisterMethod(RandomMethods, "random", randomRandom)
	RegisterMethod(RandomMethods, "range", randomRange)

	RegisterProperty(RandomProperties, "seed", randomSeedProperty)
}

// randomRandom returns a uniform pseudo-random real number in the range (0, 1).
func randomRandom(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
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

// randomSeed sets the referenced number as the seed for the pseudo-random
// generator used by the random module. If no value is passed, the current unix
// nano timestamp will be used.
func randomSeed(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if len(args) == 1 && args[0].Type() == object.NUMBER {
		seed = args[0].(*object.Number).Value.IntPart()
	} else {
		seed = time.Now().UnixNano()
	}

	rand.Seed(seed)

	return nil
}

// When called with a single number value (a), a
// pseudo-random number will be returned in the range (0, a). When called with
// two numbers (a, b), a pseudo-random number will be returned in the
// range (a, b).
func randomRange(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
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

// Properties

// randomSeedProperty returns the current seed value used internally.
func randomSeedProperty(scope *object.Scope, tok token.Token) object.Object {
	return &object.Number{Value: decimal.NewFromInt(seed)}
}
