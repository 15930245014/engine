package variable

import (
	"context"
	"errors"
	"fmt"
	com "gitlab.shudieds.com/mec/lib/entry/engine"
	"gitlab.shudieds.com/mec/lib/utils/convert"
	"gitlab.shudieds.com/mec/lib/utils/uuid"
	"gitlab.shudieds.com/zxh/engine/conf"
	_interface "gitlab.shudieds.com/zxh/engine/interface"
	"gitlab.shudieds.com/zxh/engine/types"
	eUtils "gitlab.shudieds.com/zxh/engine/utils"
	"gitlab.shudieds.com/zxh/engine/variable/adapter"
	"time"
)

type VariableCalculate struct {
	c              context.Context
	EName          string
	UniqueId       string
	params         map[string]interface{}
	variable       *com.Variable //当前的变量
	exprCalculate  _interface.ExprCalculate
	cache          map[string]types.ExprAST
	isRoot         bool
	color          map[string]bool
	variablesMp    map[string]*com.Variable //注册的变量
	factoryAdapter map[string]adapter.BaseAdapter
}

func NewVariableCalculate(c context.Context) *VariableCalculate {
	varBl := new(VariableCalculate)
	varBl.c = c
	varBl.cache = make(map[string]types.ExprAST)
	varBl.color = make(map[string]bool)
	varBl.variablesMp = make(map[string]*com.Variable)
	varBl.factoryAdapter = make(map[string]adapter.BaseAdapter)
	//
	return varBl
}
func (v *VariableCalculate) SetC(c context.Context) {
	v.c = c
}

func (v *VariableCalculate) RegisterVariables(variables []*com.Variable) {
	if len(variables) == 0 {
		return
	}
	v.variablesMp = make(map[string]*com.Variable)
	for _, variable := range variables {
		if variable.CalculateType == conf.CalculateExpr && len(variable.UniqueId) == 0 {
			variable.UniqueId = uuid.GenUuidV4()
		}
		v.variablesMp[variable.EName] = variable
	}
}
func (v *VariableCalculate) RegisterVariable(variable *com.Variable) {
	if variable == nil {
		return
	}
	v.variablesMp[variable.EName] = variable
}

func (v *VariableCalculate) Set(eName string, params map[string]interface{}, isRoot bool, obj *com.CalculateObj) error {
	v.EName = eName
	v.params = params
	v.isRoot = isRoot
	if isRoot == true {
		if _, ok := v.cache[eName]; ok {
			return nil
		}
		variable, bl := v.variablesMp[eName]
		if bl == false {
			return errors.New(eName + " load err from variables map")
		}
		v.UniqueId = variable.UniqueId
		v.variable = &com.Variable{
			Namespace:    variable.Namespace,
			EName:        variable.EName,
			CalculateObj: variable.CalculateObj,
		}
	} else {
		v.variable = &com.Variable{
			CalculateObj: *obj,
		}
	}
	return nil
}

func (v *VariableCalculate) ParseVariable() (types.ExprAST, error) {
	eName := v.EName
	returnType := v.variable.ReturnType
	isRoot := v.isRoot
	arrayType := v.variable.ArrayType
	calculateType := v.variable.CalculateType
	dftVal := v.variable.DftVal
	if isRoot {
		if _, ok := v.cache[eName]; ok {
			v.color[eName] = false
			return v.cache[eName], nil
		}
	}
	if v.color[eName] == true {
		return nil, errors.New(eName + "parse cycle err")
	}
	//if isRoot && v.color[eName] == true {
	//	return nil, errors.New(eName + "parse cycle err")
	//}
	//if v.color[eName] == true {
	//	return nil, nil
	//}
	if calculateType == conf.CalculateExpr {
		v.color[eName] = true
	}

	Adapter, err := v.GetInstanceAdapter(calculateType)
	if err != nil {
		return nil, err
	}

	Adapter.Set(v.variable, v.params, v.exprCalculate, v)
	val, err := Adapter.GetValue()
	v.color[eName] = false
	if err != nil {
		//变量设置默认值
		if len(dftVal) > 0 {
			var1 := &com.Variable{
				Namespace: "",
				EName:     eName + "-" + uuid.GenUuidV4(),
				CalculateObj: com.CalculateObj{
					CalculateType: conf.CalculateExpr,
					ReturnType:    returnType,
					Expr:          dftVal,
				},
			}
			v.RegisterVariable(var1)
			err = v.Set(var1.EName, v.params, true, nil)
			if err != nil {
				return nil, err
			}
			val2, err := v.ParseVariable()
			//没有设置默认值
			if err != nil {
				return nil, errors.New(eName + " 变量兜底默认值解析失败：" + err.Error())
			}
			val = val2.GetVal()
		} else {
			return nil, err
		}
	}
	if eUtils.IsNil(val) {
		return nil, errors.New(fmt.Sprintf("%s  exec err: value is nil", eName))
	}
	var expr types.ExprAST
	switch returnType {
	case conf.ReturnInt: //int
		intVal, ok := convert.Int(val)
		if ok == false {
			return nil, errors.New(fmt.Sprintf("%s  exec err: value is not int!", eName))
		}
		expr = types.IntExprAST{
			Val:  intVal,
			Name: eName,
		}
	case conf.ReturnFloat: //int && float64
		float64Val, ok := convert.Float64(val)
		if ok == false {
			intVal, ok := convert.Int(val)
			if ok == false {
				return nil, errors.New(fmt.Sprintf("%s  exec err: value is not float!", eName))
			}
			float64Val = float64(intVal)
		}
		expr = types.FloatExprAST{
			Val:  float64Val,
			Name: eName,
		}
	case conf.ReturnStr:
		strVal, ok := convert.String(val)
		if ok == false {
			return nil, errors.New(fmt.Sprintf("%s exec err: value is not string!", eName))
		}
		expr = types.StrExprAST{
			Val:  strVal,
			Name: eName,
		}
	case conf.ReturnBool:
		blVal, ok := convert.Bool(val)
		if ok == false {
			return nil, errors.New(fmt.Sprintf("%s  exec err: value is not bool!", eName))
		}
		expr = types.BoolExprAST{
			Val:  blVal,
			Name: eName,
		}

	case conf.ReturnMap:
		//断言
		objVal, ok := convert.ToMap(val)
		if ok == false {
			return nil, errors.New(fmt.Sprintf("%s  exec err: value is not obj!", eName))
		}
		expr = types.MapExprAST{
			Name: eName,
			Val:  objVal,
		}
	case conf.ReturnArray:
		_, ok := eUtils.ArrType(val, 1)
		if ok == false {
			return nil, errors.New(fmt.Sprintf("%s  exec err: value is not arr!", eName))
		}
		switch arrayType {
		case conf.ReturnInt:
			expr = types.ArrayExprAST{
				Name:   eName,
				IntVal: val.([]int),
				Type:   arrayType,
			}
		case conf.ReturnFloat:
			expr = types.ArrayExprAST{
				Name:     eName,
				FloatVal: val.([]float64),
				Type:     arrayType,
			}
		case conf.ReturnStr:
			expr = types.ArrayExprAST{
				Name:   eName,
				StrVal: val.([]string),
				Type:   arrayType,
			}
		case conf.ReturnBool:
			expr = types.ArrayExprAST{
				Name:    eName,
				BoolVal: val.([]bool),
				Type:    arrayType,
			}
		case conf.ReturnMap:
			//兼容从外部引用
			vv, bl := convert.ToArrMap(val)
			if bl == false {
				if ok == false {
					return nil, errors.New(fmt.Sprintf("%s  exec err: value is not []map!", eName))
				}
			}
			expr = types.ArrayExprAST{
				Name:   eName,
				MapVal: vv,
				Type:   arrayType,
			}

		case conf.ReturnDate:
			expr = types.ArrayExprAST{
				Name:    eName,
				DateVal: val.([]time.Time),
				Type:    arrayType,
			}

		case conf.ReturnAny:
			expr = types.ArrayExprAST{
				Name:   eName,
				AnyVal: val.([]interface{}),
				Type:   arrayType,
			}

		case conf.ReturnArray:
			expr = types.ArrayExprAST{
				Name:   eName,
				ArrVal: val.([]interface{}),
				Type:   arrayType,
			}
		}

	case conf.ReturnDate:
		dateVal, ok := convert.Date(val)
		if ok == false {
			return nil, errors.New(fmt.Sprintf("%s  exec err: value is not date!", eName))
		}
		expr = types.DateExprAST{
			Val:  dateVal,
			Name: eName,
		}
	case conf.ReturnAny:
		expr = types.AnyExprAST{
			Name: eName,
			Val:  val,
		}
	case conf.ReturnDecimal:
		decimalVal, ok := convert.Decimal(val)
		if ok == false {
			return nil, errors.New(fmt.Sprintf("%s  exec err: value is not decimal!", eName))
		}

		expr = types.DecimalExprAST{
			Name: eName,
			Val:  decimalVal,
		}
	default:
		return nil, errors.New(fmt.Sprintf("%s  exec err: returnType is undefined!!", eName))
	}
	//if isRoot {
	//	v.cache[eName] = expr
	//}
	v.cache[eName] = expr
	return expr, nil
}

/**
  设置公式
*/

func (v *VariableCalculate) SetExprCalculate(expression _interface.ExprCalculate) {
	v.exprCalculate = expression
}

func (v *VariableCalculate) GetInstanceAdapter(calculateType string) (adapter.BaseAdapter, error) {
	if _, ok := v.factoryAdapter[calculateType]; ok {
		return v.factoryAdapter[calculateType], nil
	}

	switch calculateType {
	case conf.CalculateInput:
		v.factoryAdapter[calculateType] = &adapter.InputAdapter{}
		break
	case conf.CalculateFixed:
		v.factoryAdapter[calculateType] = &adapter.FixedAdapter{}
		break
	case conf.CalculateExpr:
		v.factoryAdapter[calculateType] = &adapter.ExprAdapter{}
		break
	case conf.CalculateMap:
		v.factoryAdapter[calculateType] = &adapter.MapAdapter{}
		break
	case conf.CalculateArray:
		v.factoryAdapter[calculateType] = &adapter.ArrayAdapter{}
		break
	default:
		return nil, errors.New("calculate adapter is undefined!")
	}
	return v.factoryAdapter[calculateType], nil

}

func (v *VariableCalculate) ClearAllCache(ignore map[string]bool) {
	if len(ignore) != 0 {
		for k, _ := range v.cache {
			if !ignore[k] {
				delete(v.cache, k)
			}
		}

	} else {
		v.cache = make(map[string]types.ExprAST)
	}
	v.color = make(map[string]bool)
	v.factoryAdapter = make(map[string]adapter.BaseAdapter)
}
func (v *VariableCalculate) ClearVariableCache(eName string) {
	if _, ok := v.cache[eName]; ok {
		delete(v.cache, eName)
	}
}
