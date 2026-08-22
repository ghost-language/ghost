package modules

import (
	"math/rand"
	"time"

	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

var seed int64
var randomizer *rand.Rand

var RandomMethods = map[string]*object.LibraryFunction{}
var RandomProperties = map[string]*object.LibraryProperty{}

func init() {
	seed = time.Now().UnixNano()
	randomizer = rand.New(rand.NewSource(seed))

	RegisterMethod(RandomMethods, "seed", randomSeed)
	RegisterMethod(RandomMethods, "random", randomRandom)

	RegisterProperty(RandomProperties, "seed", randomSeedProperty)
}

// randomRandom returns a uniform pseudo-random real number in the range (0, 1).
func randomRandom(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	min := float64(0)
	max := float64(1)

	if len(args) > 0 {
		max = args[0].(*object.Number).Float64()

		if len(args) > 1 {
			min = max
			max = args[1].(*object.Number).Float64()
		}
	}

	number := float64(0)

	if max > 0 {
		number = float64(min + randomizer.Float64()*(max-min))
	} else {
		number = randomizer.Float64()
	}

	return object.NewFloat(number)
}

// randomSeed sets the referenced number as the seed for the pseudo-random
// generator used by the random module. If no value is passed, the current unix
// nano timestamp will be used.
func randomSeed(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if len(args) == 1 && args[0].Type() == object.NUMBER {
		SeedRandom(args[0].(*object.Number).Int64())
	} else {
		SeedRandom(time.Now().UnixNano())
	}

	return nil
}

// SeedRandom points the generator behind the random and math modules at a given
// seed. It is exported so that a host embedding Ghost can fix the seed before a
// program runs - an engine that wants its procedural output reproducible by
// default, say - without reaching into the module tables to do it.
func SeedRandom(value int64) {
	seed = value
	randomizer = rand.New(rand.NewSource(seed))
}

// Properties

// randomSeedProperty returns the current seed value used internally.
func randomSeedProperty(scope *object.Scope, tok token.Token) object.Object {
	return object.NewInt(seed)
}
