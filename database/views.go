package views

import (
	"github.com/kr/pretty"
	"github.com/robertkrimen/otto"
)

func MapFun(js string, key string, val map[string]interface{}) {
	vm := otto.New()
	vm.Set("emit", func(call otto.FunctionCall) otto.Value {
		key := call.Argument(0)
		value := call.Argument(1)

		// indexing happens here
		pretty.Log("emitting", key, value)

		return otto.Value{}
	})
	_, err := vm.Call(js, nil, map[string]interface{}{
		"_id":  key,
		"_val": val,
	})

	if err != nil {
		pretty.Log("mapfun error ", err)
	}
}
