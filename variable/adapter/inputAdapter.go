package adapter

import (
	"errors"
	"fmt"
	com "gitlab.shudieds.com/mec/lib/entry/engine"
	_interface "gitlab.shudieds.com/zxh/engine/interface"
	"gitlab.shudieds.com/zxh/engine/utils"
)

type InputAdapter struct {
	variable *com.Variable
	params   map[string]interface{}
}

func (adapter *InputAdapter) Set(variable *com.Variable, params map[string]interface{}, expression _interface.ExprCalculate, variable2 _interface.Variable) {
	adapter.variable = variable
	adapter.params = params
}

func (adapter *InputAdapter) GetValue() (interface{}, error) {
	var val = adapter.variable.Val
	_, ok := adapter.params[adapter.variable.EName]

	if !ok && utils.IsNil(val) {
		return nil, errors.New(fmt.Sprintf("%s inputAdapter exec err: params not exists key", adapter.variable.EName))
	}
	if ok {
		val = adapter.params[adapter.variable.EName]
	}

	return val, nil

}
