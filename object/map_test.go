package object

import "testing"

// keysOf reads the keys back out of pairs, for comparing against an
// expected order without depending on any one Object type's String().
func keysOf(pairs []MapPair) []string {
	keys := make([]string, len(pairs))

	for index, pair := range pairs {
		keys[index] = pair.Key.(*String).Value
	}

	return keys
}

func assertKeyOrder(t *testing.T, pairs []MapPair, expected ...string) {
	t.Helper()

	got := keysOf(pairs)

	if len(got) != len(expected) {
		t.Fatalf("wrong number of pairs. got=%v, expected=%v", got, expected)
	}

	for index, key := range expected {
		if got[index] != key {
			t.Fatalf("wrong order. got=%v, expected=%v", got, expected)
		}
	}
}

func setString(mapObject *Map, key string, value int64) {
	pairKey := &String{Value: key}
	mapObject.SetPair(pairKey.MapKey(), MapPair{Key: pairKey, Value: NewInt(value)})
}

// TestMapSetPairPreservesInsertionOrder covers §13.5's core mechanism: a new
// key is appended, and an existing key keeps its original position - only
// its value changes.
func TestMapSetPairPreservesInsertionOrder(t *testing.T) {
	mapObject := NewOrderedMap()

	setString(mapObject, "z", 1)
	setString(mapObject, "a", 2)
	setString(mapObject, "m", 3)

	assertKeyOrder(t, mapObject.OrderedPairs(), "z", "a", "m")

	// Re-setting "a" must not move it to the end.
	setString(mapObject, "a", 99)

	assertKeyOrder(t, mapObject.OrderedPairs(), "z", "a", "m")

	pair, ok := mapObject.Pairs[(&String{Value: "a"}).MapKey()]

	if !ok || pair.Value.(*Number).Int64() != 99 {
		t.Fatalf("re-setting an existing key should update its value. got=%+v", pair)
	}
}

// TestMapRemovePairKeepsOrderConsistent confirms a removed key leaves no gap,
// and a key added afterward is appended at the end as usual.
func TestMapRemovePairKeepsOrderConsistent(t *testing.T) {
	mapObject := NewOrderedMap()

	setString(mapObject, "a", 1)
	setString(mapObject, "b", 2)
	setString(mapObject, "c", 3)

	removed, ok := mapObject.RemovePair((&String{Value: "b"}).MapKey())

	if !ok || removed.Value.(*Number).Int64() != 2 {
		t.Fatalf("RemovePair should answer the pair that was there. got=%+v, ok=%v", removed, ok)
	}

	assertKeyOrder(t, mapObject.OrderedPairs(), "a", "c")

	setString(mapObject, "d", 4)

	assertKeyOrder(t, mapObject.OrderedPairs(), "a", "c", "d")

	// Removing a key that was never present changes nothing.
	_, ok = mapObject.RemovePair((&String{Value: "missing"}).MapKey())

	if ok {
		t.Fatal("RemovePair should answer false for a key that was never present")
	}

	assertKeyOrder(t, mapObject.OrderedPairs(), "a", "c", "d")
}

// TestMapOrderedPairsIsStable confirms repeated reads of the same Map answer
// the identical order every time - the property a bare Go map's randomized
// iteration cannot guarantee even within a single run.
func TestMapOrderedPairsIsStable(t *testing.T) {
	mapObject := NewOrderedMap()

	setString(mapObject, "z", 1)
	setString(mapObject, "a", 2)
	setString(mapObject, "m", 3)
	setString(mapObject, "b", 4)

	first := keysOf(mapObject.OrderedPairs())

	for i := 0; i < 50; i++ {
		assertKeyOrder(t, mapObject.OrderedPairs(), first...)
	}
}
