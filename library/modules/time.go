package modules

import (
	"time"

	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
	"github.com/shopspring/decimal"
)

var TimeMethods = map[string]*object.LibraryFunction{}
var TimeProperties = map[string]*object.LibraryProperty{}

func init() {
	RegisterMethod(TimeMethods, "sleep", timeSleep)
	RegisterMethod(TimeMethods, "now", timeNow)

	RegisterProperty(TimeProperties, "nanosecond", timeNanosecond)
	RegisterProperty(TimeProperties, "microsecond", timeMicrosecond)
	RegisterProperty(TimeProperties, "millisecond", timeMillisecond)
	RegisterProperty(TimeProperties, "second", timeSecond)
	RegisterProperty(TimeProperties, "minute", timeMinute)
	RegisterProperty(TimeProperties, "hour", timeHour)
	RegisterProperty(TimeProperties, "day", timeDay)
	RegisterProperty(TimeProperties, "week", timeWeek)
	RegisterProperty(TimeProperties, "month", timeMonth)
	RegisterProperty(TimeProperties, "year", timeYear)
}

func timeSleep(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 1 {
		// TODO: error
		return nil
	}

	if args[0].Type() != object.NUMBER {
		// TODO: error
		return nil
	}

	ms := args[0].(*object.Number)
	time.Sleep(time.Duration(ms.Value.IntPart()) * time.Millisecond)

	return nil
}

func timeNow(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 0 {
		// TODO: error
		return nil
	}

	unix := decimal.NewFromInt(time.Now().Unix())

	return &object.Number{Value: unix}
}

// properties

func timeNanosecond(scope *object.Scope, tok token.Token) object.Object {
	nanosecond := decimal.NewFromFloat(0.00001)

	return &object.Number{Value: nanosecond}
}

func timeMicrosecond(scope *object.Scope, tok token.Token) object.Object {
	microsecond := decimal.NewFromFloat(0.0001)

	return &object.Number{Value: microsecond}
}

func timeMillisecond(scope *object.Scope, tok token.Token) object.Object {
	millisecond := decimal.NewFromFloat(0.001)

	return &object.Number{Value: millisecond}
}

func timeSecond(scope *object.Scope, tok token.Token) object.Object {
	second := decimal.NewFromInt(1)

	return &object.Number{Value: second}
}

func timeMinute(scope *object.Scope, tok token.Token) object.Object {
	minute := decimal.NewFromInt(60)

	return &object.Number{Value: minute}
}

func timeHour(scope *object.Scope, tok token.Token) object.Object {
	hour := decimal.NewFromInt(3600)

	return &object.Number{Value: hour}
}

func timeDay(scope *object.Scope, tok token.Token) object.Object {
	day := decimal.NewFromInt(86400)

	return &object.Number{Value: day}
}

func timeWeek(scope *object.Scope, tok token.Token) object.Object {
	week := decimal.NewFromInt(604800)

	return &object.Number{Value: week}
}

func timeMonth(scope *object.Scope, tok token.Token) object.Object {
	month := decimal.NewFromInt(2592000)

	return &object.Number{Value: month}
}

func timeYear(scope *object.Scope, tok token.Token) object.Object {
	year := decimal.NewFromInt(31536000)

	return &object.Number{Value: year}
}
