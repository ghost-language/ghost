package object

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
	"github.com/shopspring/decimal"
)

var evaluator func(node ast.Node, scope *Scope) Object

// Type is the type of the object given as a string.
type Type string

// Object is the interface for all object values.
type Object interface {
	HasMethods
	Type() Type
	String() string
}

type MapKey struct {
	Type  Type
	Value uint64
}

type Mappable interface {
	MapKey() MapKey
}

type HasMethods interface {
	Method(method string, args []Object) (Object, bool)
}

type GoFunction func(scope *Scope, tok token.Token, args ...Object) Object
type GoProperty func(scope *Scope, tok token.Token) Object
type ObjectMethod func(value interface{}, args ...Object) (Object, bool)

func RegisterEvaluator(e func(node ast.Node, scope *Scope) Object) {
	evaluator = e
}

func AnyValueToObject(val any) Object {
	switch v := val.(type) {
	case bool:
		if v {
			return &Boolean{Value: true}
		}

		return &Boolean{Value: false}
	case string:
		return &String{Value: v}
	case int:
		return &Number{Value: decimal.NewFromInt(int64(v))}
	case int64:
		return &Number{Value: decimal.NewFromInt(int64(v))}
	case float64:
		return &Number{Value: decimal.NewFromFloat(v)}
	case nil:
		return &Null{}
	case []any:
		elements := make([]Object, len(v))

		for index, item := range v {
			elements[index] = AnyValueToObject(item)
		}

		return &List{Elements: elements}
	case map[string]any:
		pairs := make(map[MapKey]MapPair)

		for key, val := range v {
			pairKey := &String{Value: key}
			var pairValue Object
			hashed := pairKey.MapKey()

			pairValue = AnyValueToObject(val)

			pairs[hashed] = MapPair{Key: pairKey, Value: pairValue}
		}

		return &Map{Pairs: pairs}
	}

	return nil
}

func ObjectToAnyValue(val Object) any {
	switch v := val.(type) {
	case *Boolean:
		return bool(v.Value)
	case *String:
		return string(v.Value)
	case *Number:
		// Determine if value is an integer or float.
		if v.Value.Exponent() <= 0 {
			return int(v.Value.IntPart())
		}

		num, _ := v.Value.Float64()

		return num
	case *Null:
		return nil
	case *List:
		var collection []any

		for _, val := range v.Elements {
			collection = append(collection, ObjectToAnyValue(val))
		}

		return collection
	case *Map:
		collection := make(map[string]any)

		for _, pair := range v.Pairs {
			collection[pair.Key.(*String).Value] = ObjectToAnyValue(pair.Value)
		}

		return collection
	}

	return nil
}
