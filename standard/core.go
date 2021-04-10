package standard

import (
	"io/ioutil"
	"time"

	"ghostlang.org/x/ghost/object"
	"github.com/shopspring/decimal"
)

func init() {
	RegisterFunction("clock", clockFunction)
	RegisterFunction("file", fileFunction)
}

func clockFunction(args []object.Object) object.Object {
	secs := decimal.NewFromInt(time.Now().Unix())

	return &object.Number{Value: secs}
}

func fileFunction(args []object.Object) object.Object {
	path := args[0].String()

	contents, err := ioutil.ReadFile(path)

	if err != nil {
		return &object.Error{Message: err.Error()}
	}

	return &object.String{Value: string(contents)}
}