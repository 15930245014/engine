package adapter

import (
	"errors"
	"fmt"
	com "gitlab.shudieds.com/mec/lib/entry/engine"
	_interface "gitlab.shudieds.com/zxh/engine/interface"
	eUtils "gitlab.shudieds.com/zxh/engine/utils"
)

type InputAdapter struct {
	field  *com.Field
	params map[string]interface{}
}

func (adapter *InputAdapter) Set(field *com.Field, params map[string]interface{}, expression _interface.ExprCalculate, f _interface.Field) {
	adapter.field = field
	adapter.params = params
}

func (adapter *InputAdapter) GetValue() (interface{}, error) {
	var val = adapter.field.Val
	_, ok := adapter.params[adapter.field.EName]

	if !ok && eUtils.IsNil(val) {
		return nil, errors.New(fmt.Sprintf("%s inputAdapter exec err: params not exists key", adapter.field.EName))
	}
	if ok {
		val = adapter.params[adapter.field.EName]
	}

	return val, nil

}
