package adapter

import (
	"errors"
	"fmt"
	com "gitlab.shudieds.com/mec/lib/entry/engine"
	"gitlab.shudieds.com/mec/lib/utils/jsoniter"
	"gitlab.shudieds.com/mec/lib/utils/uuid"
	"gitlab.shudieds.com/zxh/engine/conf"
	_interface "gitlab.shudieds.com/zxh/engine/interface"
	eUtils "gitlab.shudieds.com/zxh/engine/utils"
)

type MapAdapter struct {
	field      *com.Field
	params     map[string]interface{}
	fieldParse _interface.Field
}

func (adapter *MapAdapter) Set(field *com.Field, params map[string]interface{}, expression _interface.ExprCalculate, f _interface.Field) {
	adapter.fieldParse = f
	adapter.params = params
	adapter.field = field
}

func (adapter *MapAdapter) GetValue() (interface{}, error) {
	//解析obj
	if eUtils.IsNil(adapter.field.Map) {
		return nil, errors.New(fmt.Sprintf("%s MapAdapter GetValue obj is nil", adapter.field.EName))
	}
	fName := adapter.field.EName
	objV := make(map[string]com.CalculateObj)
	for k, v := range adapter.field.Map {
		cbj := com.CalculateObj{
			CalculateType: v.CalculateType,
			ReturnType:    v.ReturnType,
			Val:           v.Val,
			Expr:          v.Expr,
			ArrayType:     v.ArrayType,
			Map:           v.Map,
		}
		if len(v.Array) > 0 {
			t := make([]com.CalculateObj, len(v.Array))
			copy(t, v.Array)
			cbj.Array = t
		}
		if len(v.Map) > 0 {
			bt, _ := jsoniter.Marshal(v.Map)
			err := jsoniter.Unmarshal(bt, &cbj.Map)
			if err != nil {
				return nil, errors.New(fName + ":map配置格式错误")
			}
		}
		if len(cbj.UniqueId) == 0 && v.CalculateType == conf.CalculateExpr {
			cbj.UniqueId = uuid.GenUuidV3(v.Expr)
		}
		objV[k] = cbj

	}
	rt := make(map[string]interface{})
	for k, val := range objV {
		err := adapter.fieldParse.Set(fName+"-"+k, adapter.params, false, &val)
		if err != nil {
			return nil, err
		}
		fv, err := adapter.fieldParse.ParseField()

		if err != nil {
			return nil, err
		}
		if fv != nil { //字段内部引用
			rt[k] = fv.GetVal()
		}
	}
	return rt, nil
}
