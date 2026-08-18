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
	variable      *com.Variable
	params        map[string]interface{}
	variableParse _interface.Variable
}

func (adapter *MapAdapter) Set(variable *com.Variable, params map[string]interface{}, expression _interface.ExprCalculate, variable2 _interface.Variable) {
	adapter.variable = variable
	adapter.params = params
	adapter.variableParse = variable2
}

func (adapter *MapAdapter) GetValue() (interface{}, error) {
	//解析obj
	if eUtils.IsNil(adapter.variable.Map) {
		return nil, errors.New(fmt.Sprintf("%s ObjAdapter GetValue obj is nil", adapter.variable.EName))
	}
	EName := adapter.variable.EName
	//深拷贝
	objV := make(map[string]com.CalculateObj)
	for k, v := range adapter.variable.Map {
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
		//彻底解决map深拷贝问题
		if len(v.Map) > 0 {
			bt, _ := jsoniter.Marshal(v.Map)
			err := jsoniter.Unmarshal(bt, &cbj.Map)
			if err != nil {
				return nil, errors.New(EName + ":map配置格式错误")
			}
		}
		if len(cbj.UniqueId) == 0 && v.CalculateType == conf.CalculateExpr {
			cbj.UniqueId = uuid.GenUuidV3(v.Expr)
		}
		objV[k] = cbj
	}

	rt := make(map[string]interface{})
	for k, val := range objV {
		err := adapter.variableParse.Set(EName+"-"+k, adapter.params, false, &val)

		if err != nil {
			return nil, err
		}
		fv, err := adapter.variableParse.ParseVariable()

		if err != nil {
			return nil, err
		}
		if fv != nil {
			rt[k] = fv.GetVal()
		}
	}
	return rt, nil
}
