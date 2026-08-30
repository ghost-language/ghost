package modules

import (
	"sync"
	"testing"

	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

// callRandom invokes a registered random method the way the evaluator would.
func callRandom(t *testing.T, name string, args ...object.Object) object.Object {
	t.Helper()

	method, ok := RandomMethods[name]

	if !ok {
		t.Fatalf("random.%s is not registered", name)
	}

	return method.Function(nil, token.Token{}, args...)
}

func floatOf(t *testing.T, result object.Object) float64 {
	t.Helper()

	number, ok := result.(*object.Number)

	if !ok {
		t.Fatalf("object is not Number. got=%T (%+v)", result, result)
	}

	return number.Float64()
}

func TestRandomRandomRange(t *testing.T) {
	for i := 0; i < 50; i++ {
		value := floatOf(t, callRandom(t, "random"))

		if value < 0 || value >= 1 {
			t.Fatalf("random() out of [0, 1): got=%v", value)
		}
	}

	for i := 0; i < 50; i++ {
		value := floatOf(t, callRandom(t, "random", object.NewFloat(10)))

		if value < 0 || value >= 10 {
			t.Fatalf("random(10) out of [0, 10): got=%v", value)
		}
	}

	for i := 0; i < 50; i++ {
		value := floatOf(t, callRandom(t, "random", object.NewFloat(5), object.NewFloat(10)))

		if value < 5 || value >= 10 {
			t.Fatalf("random(5, 10) out of [5, 10): got=%v", value)
		}
	}
}

// TestRandomSeedIsReproducible covers the same seeding contract
// TestMathRandomIsReproducible covers for math.randomInt: a given seed always
// replays the same sequence.
func TestRandomSeedIsReproducible(t *testing.T) {
	callRandom(t, "seed", object.NewInt(1234))
	first := floatOf(t, callRandom(t, "random"))

	callRandom(t, "seed", object.NewInt(1234))
	second := floatOf(t, callRandom(t, "random"))

	if first != second {
		t.Errorf("seeded runs differ. got=%v and %v", first, second)
	}
}

func TestRandomCurrentSeedProperty(t *testing.T) {
	callRandom(t, "seed", object.NewInt(99))

	property, ok := RandomProperties["currentSeed"]

	if !ok {
		t.Fatal("random.currentSeed is not registered")
	}

	got := property.Property(nil, token.Token{})

	number, ok := got.(*object.Number)

	if !ok {
		t.Fatalf("currentSeed is not Number. got=%T", got)
	}

	if number.Int64() != 99 {
		t.Errorf("currentSeed: got=%d, expected=99", number.Int64())
	}
}

// TestRandomConcurrentAccessIsRaceFree covers §13.1: random.random() and
// math.randomInt() share one *rand.Rand behind randomState, and *rand.Rand is
// not safe for concurrent use on its own. Before that state was guarded, this
// test failed under `go test -race` (and could panic or corrupt the sequence
// even without -race) - it is the regression test for that fix. Ghost itself
// reaches this concurrently through http.handle(), which runs each request's
// callback on its own goroutine.
func TestRandomConcurrentAccessIsRaceFree(t *testing.T) {
	// Read the two functions on the test's own goroutine - t.Fatal (inside
	// callRandom/call) is only safe to reach from there, and the goroutines
	// below only need the already-resolved *object.LibraryFunction.
	random, ok := RandomMethods["random"]

	if !ok {
		t.Fatal("random.random is not registered")
	}

	randomInt, ok := MathMethods["randomInt"]

	if !ok {
		t.Fatal("math.randomInt is not registered")
	}

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(2)

		go func() {
			defer wg.Done()

			random.Function(nil, token.Token{})
		}()

		go func() {
			defer wg.Done()

			randomInt.Function(nil, token.Token{}, object.NewInt(1), object.NewInt(1000))
		}()
	}

	wg.Wait()
}
