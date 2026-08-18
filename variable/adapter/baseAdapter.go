package adapter

import (
	com "gitlab.shudieds.com/mec/lib/entry/engine"
	_interface "gitlab.shudieds.com/zxh/engine/interface"
)

type BaseAdapter interface {
	Set(*com.Variable, map[string]interface{}, _interface.ExprCalculate, _interface.Variable)
	GetValue() (interface{}, error)
}
