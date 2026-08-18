package adapter

import (
	com "gitlab.shudieds.com/mec/lib/entry/engine"
	_interface "gitlab.shudieds.com/zxh/engine/interface"
)

type FixedAdapter struct {
	variable *com.Variable
	params   map[string]interface{}
}

func (adapter *FixedAdapter) Set(variable *com.Variable, params map[string]interface{}, expression _interface.ExprCalculate, variable2 _interface.Variable) {
	adapter.variable = variable
	adapter.params = params
}

func (adapter *FixedAdapter) GetValue() (interface{}, error) {
	val := adapter.variable.Val
	return val, nil
}
