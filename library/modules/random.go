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

	RegisterProperty(RandomProperties, "currentSeed", randomCurrentSeed)
}

// randomRandom returns a uniform pseudo-random real number in the range (0, 1).
func randomRandom(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arityRange("random.random", tok, args, 0, 2); err != nil {
		return err
	}

	min := float64(0)
	max := float64(1)

	if len(args) > 0 {
		bound, err := floatAt("random.random", tok, args, 0)

		if err != nil {
			return err
		}

		max = bound

		if len(args) > 1 {
			upper, err := floatAt("random.random", tok, args, 1)

			if err != nil {
				return err
			}

			min = max
			max = upper
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
	if err := arityRange("random.seed", tok, args, 0, 1); err != nil {
		return err
	}

	if len(args) == 0 {
		SeedRandom(time.Now().UnixNano())

		return nil
	}

	chosen, err := integerAt("random.seed", tok, args, 0)

	if err != nil {
		return err
	}

	SeedRandom(chosen)

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

// randomCurrentSeed returns the seed currently driving the generator. It is a
// property rather than a method, and named apart from seed() - the method
// that sets it - so that setting and reading the seed are not two different
// calls sharing one name.
func randomCurrentSeed(scope *object.Scope, tok token.Token) object.Object {
	return object.NewInt(seed)
}
