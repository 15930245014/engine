package adapter

import (
	"errors"
	"fmt"
	com "gitlab.shudieds.com/mec/lib/entry/engine"
	"gitlab.shudieds.com/mec/lib/utils/maphelper"
	_interface "gitlab.shudieds.com/zxh/engine/interface"
	"strings"
)

type FuncAdapter struct {
	field      *com.Field
	params     map[string]interface{}
	fieldParse _interface.Field
}

func (adapter *FuncAdapter) Set(field *com.Field, params map[string]interface{}, expression _interface.ExprCalculate, f _interface.Field) {
	adapter.fieldParse = f
	adapter.params = params
	adapter.field = field
}

func (adapter *FuncAdapter) GetValue() (interface{}, error) {
	fieldArr := strings.Split(adapter.field.EName, ".")
	if len(fieldArr) < 2 {
		return nil, errors.New(fmt.Sprintf("%s funcAdapter 字段名称不合法 is nil", adapter.field.EName))
	}
	err := adapter.fieldParse.Set(fieldArr[0], adapter.params, true, nil)
	if err != nil {
		return nil, err
	}
	fv, err := adapter.fieldParse.ParseField()
	if err != nil {
		return nil, err
	}
	rt, bl := fv.GetVal().(map[string]interface{})
	if bl == false {
		return nil, errors.New(fmt.Sprintf("%s funcAdapter 返回值类型不合法 is nil", adapter.field.EName))
	}

	val, bl := maphelper.JDGet(rt, fieldArr[1:])
	if bl == false {
		return nil, errors.New(fmt.Sprintf("%s funcAdapter 返回值结构和取值不一致", adapter.field.EName))
	}

	return val, nil
}
