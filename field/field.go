package field

import (
	"context"
	"errors"
	"fmt"
	com "gitlab.shudieds.com/mec/lib/entry/engine"
	"gitlab.shudieds.com/mec/lib/utils/convert"
	"gitlab.shudieds.com/mec/lib/utils/uuid"
	"gitlab.shudieds.com/zxh/engine/conf"
	"gitlab.shudieds.com/zxh/engine/field/adapter"
	_interface "gitlab.shudieds.com/zxh/engine/interface"
	"gitlab.shudieds.com/zxh/engine/types"
	eUtils "gitlab.shudieds.com/zxh/engine/utils"
	"time"
)

/**
连接器字符解析
*/

type FieldCalculate struct {
	c              context.Context
	eName          string
	params         map[string]interface{}
	field          *com.Field
	exprCalculate  _interface.ExprCalculate
	cache          map[string]types.ExprAST
	isRoot         bool
	color          map[string]bool
	fieldsMp       map[string]*com.Field
	factoryAdapter map[string]adapter.BaseAdapter
}

func NewFieldCalculate(c context.Context) *FieldCalculate {
	f := new(FieldCalculate)
	f.c = c
	f.cache = make(map[string]types.ExprAST)
	f.color = make(map[string]bool)
	f.fieldsMp = make(map[string]*com.Field)
	f.factoryAdapter = make(map[string]adapter.BaseAdapter)
	return f
}
func (f *FieldCalculate) SetC(c context.Context) {
	f.c = c
}
func (v *FieldCalculate) RegisterFields(fields []*com.Field) {
	if len(fields) == 0 {
		return
	}
	v.fieldsMp = make(map[string]*com.Field)
	for _, field := range fields {
		if field.CalculateType == conf.CalculateExpr && len(field.UniqueId) == 0 {
			field.UniqueId = uuid.GenUuidV4()
		}
		v.fieldsMp[field.EName] = field
	}
}
func (v *FieldCalculate) RegisterField(field *com.Field) {
	if field == nil {
		return
	}
	v.fieldsMp[field.EName] = field
}

func (v *FieldCalculate) Set(eName string, params map[string]interface{}, isRoot bool, obj *com.CalculateObj) error {
	v.eName = eName
	v.isRoot = isRoot
	v.params = params
	if isRoot {
		if _, ok := v.cache[eName]; ok {
			return nil
		}
		field, bl := v.fieldsMp[eName]
		if bl == false {
			return errors.New(eName + " load err from fields map")
		}
		v.field = &com.Field{
			Namespace:    field.Namespace,
			EName:        field.EName,
			CalculateObj: field.CalculateObj,
		}
	} else {
		v.field = &com.Field{
			CalculateObj: *obj,
			EName:        eName,
		}
	}
	return nil
}
func (v *FieldCalculate) ParseField() (types.ExprAST, error) {
	EName := v.eName
	returnType := v.field.ReturnType
	//isRoot := v.isRoot
	arrayType := v.field.ArrayType
	calculateType := v.field.CalculateType
	dftVal := v.field.DftVal
	if _, ok := v.cache[EName]; ok {
		v.color[EName] = false
		return v.cache[EName], nil
	}

	if v.color[EName] == true {
		return nil, errors.New(EName + " :parse cycle err")
	}
	//if isRoot == true && v.color[EName] == true {
	//	return nil, errors.New(EName + "parse cycle err")
	//}
	//if isRoot == false && v.color[EName] == true {
	//	return nil, nil
	//}

	if calculateType == conf.CalculateExpr {
		v.color[EName] = true
	}

	Adapter, err := v.getInstanceAdapter(calculateType)
	if err != nil {
		return nil, err
	}

	Adapter.Set(v.field, v.params, v.exprCalculate, v)
	val, err := Adapter.GetValue()
	v.color[EName] = false
	if err != nil {
		//变量设置默认值
		if len(dftVal) > 0 {
			var1 := &com.Field{
				Namespace: "",
				EName:     EName + "-" + uuid.GenUuidV4(),
				CalculateObj: com.CalculateObj{
					CalculateType: conf.CalculateExpr,
					ReturnType:    returnType,
					Expr:          dftVal,
				},
			}
			v.RegisterField(var1)
			err = v.Set(var1.EName, v.params, true, nil)
			if err != nil {
				return nil, err
			}
			val2, err := v.ParseField()
			//没有设置默认值
			if err != nil {
				return nil, errors.New(EName + " 字段兜底默认值解析失败：" + err.Error())
			}
			val = val2.GetVal()
		} else {
			return nil, err
		}
	}
	if eUtils.IsNil(val) {
		return nil, errors.New(fmt.Sprintf("%s  exec err: value is nil", EName))
	}
	var expr types.ExprAST
	switch returnType {
	case conf.ReturnInt:
		intVal, ok := convert.Int(val)
		if ok == false {
			return nil, errors.New(fmt.Sprintf("%s  exec err: value is not int!", EName))
		}
		expr = types.IntExprAST{
			Val:  intVal,
			Name: EName,
		}
	case conf.ReturnFloat:
		float64Val, ok := convert.Float64(val)
		if ok == false {
			intVal, ok := convert.Int(val)
			if ok == false {
				return nil, errors.New(fmt.Sprintf("%s  exec err: value is not float64!", EName))
			}
			float64Val = float64(intVal)
		}
		expr = types.FloatExprAST{
			Val:  float64Val,
			Name: EName,
		}
	case conf.ReturnStr:
		strVal, ok := convert.String(val)
		if ok == false {
			return nil, errors.New(fmt.Sprintf("%s exec err: value is not string!", EName))
		}
		expr = types.StrExprAST{
			Val:  strVal,
			Name: EName,
		}
	case conf.ReturnBool:
		blVal, ok := convert.Bool(val)
		if ok == false {
			return nil, errors.New(fmt.Sprintf("%s  exec err: value is not bool!", EName))
		}
		expr = types.BoolExprAST{
			Val:  blVal,
			Name: EName,
		}
	case conf.ReturnAny:
		expr = types.AnyExprAST{
			Name: EName,
			Val:  val,
		}
	case conf.ReturnMap:
		objVal, ok := convert.ToMap(val)
		if ok == false {
			return nil, errors.New(fmt.Sprintf("%s  exec err: value is not obj!", EName))
		}

		expr = types.MapExprAST{
			Name: EName,
			Val:  objVal,
		}
	case conf.ReturnArray:
		_, ok := eUtils.ArrType(val, 1)
		if ok == false {
			return nil, errors.New(fmt.Sprintf("%s  exec err: value is not arr!", EName))
		}
		switch arrayType {
		case conf.ReturnInt:
			expr = types.ArrayExprAST{
				Name:   EName,
				IntVal: val.([]int),
				Type:   arrayType,
			}
		case conf.ReturnFloat:
			expr = types.ArrayExprAST{
				Name:     EName,
				FloatVal: val.([]float64),
				Type:     arrayType,
			}
		case conf.ReturnStr:
			expr = types.ArrayExprAST{
				Name:   EName,
				StrVal: val.([]string),
				Type:   arrayType,
			}
		case conf.ReturnBool:
			expr = types.ArrayExprAST{
				Name:    EName,
				BoolVal: val.([]bool),
				Type:    arrayType,
			}
		case conf.ReturnMap:
			//兼容从外部引用
			vv, bl := convert.ToArrMap(val)
			if bl == false {
				if ok == false {
					return nil, errors.New(fmt.Sprintf("%s  exec err: value is not []map!", EName))
				}
			}
			expr = types.ArrayExprAST{
				Name:   EName,
				MapVal: vv,
				Type:   arrayType,
			}
		case conf.ReturnDate:
			expr = types.ArrayExprAST{
				Name:    EName,
				DateVal: val.([]time.Time),
				Type:    arrayType,
			}

		case conf.ReturnAny:
			expr = types.ArrayExprAST{
				Name:   EName,
				AnyVal: val.([]interface{}),
				Type:   arrayType,
			}

		case conf.ReturnArray:
			expr = types.ArrayExprAST{
				Name:   EName,
				ArrVal: val.([]interface{}),
				Type:   arrayType,
			}
		}
	case conf.ReturnDate:
		dateVal, ok := convert.Date(val)
		if ok == false {
			return nil, errors.New(fmt.Sprintf("%s  exec err: value is not date!", EName))
		}
		expr = types.DateExprAST{
			Val:  dateVal,
			Name: EName,
		}
	case conf.ReturnDecimal:
		decimalVal, ok := convert.Decimal(val)
		if ok == false {
			return nil, errors.New(fmt.Sprintf("%s  exec err: value is not deecimal!", EName))
		}

		expr = types.DecimalExprAST{
			Name: EName,
			Val:  decimalVal,
		}
	default:
		return nil, errors.New(fmt.Sprintf("%s  exec err: returnType is undefined!", EName))
	}
	v.cache[EName] = expr
	return expr, nil
}

/**
  设置公式
*/

func (v *FieldCalculate) SetExprCalculate(expression _interface.ExprCalculate) {
	v.exprCalculate = expression
}

func (v *FieldCalculate) getInstanceAdapter(calculateType string) (adapter.BaseAdapter, error) {
	if _, ok := v.factoryAdapter[calculateType]; ok {
		return v.factoryAdapter[calculateType], nil
	}
	switch calculateType {
	case conf.CalculateInput:
		v.factoryAdapter[calculateType] = &adapter.InputAdapter{}
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
	case conf.CalculateFixed:
		v.factoryAdapter[calculateType] = &adapter.FixedAdapter{}
	default:
		return nil, errors.New("Field calculate adapter is undefined!")
	}
	return v.factoryAdapter[calculateType], nil
}
func (v *FieldCalculate) ClearAllCache() {
	v.cache = make(map[string]types.ExprAST)
	v.color = make(map[string]bool)
	v.factoryAdapter = make(map[string]adapter.BaseAdapter)

}
func (v *FieldCalculate) ClearFieldCache(fName string) {
	if _, ok := v.cache[fName]; ok {
		delete(v.cache, fName)
	}
}
