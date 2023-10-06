package object

import "github.com/shopspring/decimal"

const MAP = "MAP"

// Map objects consist of a map value.
type Map struct {
	Pairs map[MapKey]MapPair
}

type MapPair struct {
	Key   Object
	Value Object
}

// String represents the map object's value as a string.
func (mapObject *Map) String() string {
	return "map"
}

// Type returns the map object type.
func (mapObject *Map) Type() Type {
	return MAP
}

// Method defines the set of methods available on map objects.
func (mapObject *Map) Method(method string, args []Object) (Object, bool) {
	return nil, false
}

func NewMap(values map[string]interface{}) *Map {
	pairs := make(map[MapKey]MapPair)

	for key, value := range values {
		pairKey := &String{Value: key}
		var pairValue Object
		hashed := pairKey.MapKey()

		switch val := value.(type) {
		case int:
		case int64:
			pairValue = &Number{Value: decimal.NewFromInt(int64(val))}
		case string:
			pairValue = &String{Value: val}
		}

		pairs[hashed] = MapPair{Key: pairKey, Value: pairValue}
	}

	return &Map{Pairs: pairs}
}
