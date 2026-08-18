package adapter

import (
	com "gitlab.shudieds.com/mec/lib/entry/engine"
	_interface "gitlab.shudieds.com/zxh/engine/interface"
)

type FixedAdapter struct {
	field  *com.Field
	params map[string]interface{}
}

func (adapter *FixedAdapter) Set(variable *com.Field, params map[string]interface{}, expression _interface.ExprCalculate, f _interface.Field) {
	adapter.field = variable
	//adapter.params = params
}

func (adapter *FixedAdapter) GetValue() (interface{}, error) {
	val := adapter.field.Val
	return val, nil
}
