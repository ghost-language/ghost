package object

const MAP = "MAP"

// Map objects consist of a map value.
type Map struct {
	Pairs map[MapKey]MapPair
}

type MapPair struct {
	Key   Object
	Value Object
}

func (object *Map) Accept(v Visitor) {
	v.visitMap(object)
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
