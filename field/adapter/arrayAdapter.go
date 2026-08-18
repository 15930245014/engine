package adapter

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	com "gitlab.shudieds.com/mec/lib/entry/engine"
	"gitlab.shudieds.com/mec/lib/utils/convert"
	"gitlab.shudieds.com/mec/lib/utils/uuid"
	"gitlab.shudieds.com/zxh/engine/conf"
	_interface "gitlab.shudieds.com/zxh/engine/interface"
	eUtils "gitlab.shudieds.com/zxh/engine/utils"
	"strconv"
	"time"
)

type ArrayAdapter struct {
	field      *com.Field
	params     map[string]interface{}
	fieldParse _interface.Field
}

func (adapter *ArrayAdapter) Set(field *com.Field, params map[string]interface{}, expression _interface.ExprCalculate, f _interface.Field) {
	adapter.field = field
	adapter.params = params
	adapter.fieldParse = f
}

func (adapter *ArrayAdapter) GetValue() (interface{}, error) {
	//解析obj
	if eUtils.IsNil(adapter.field.Array) {
		return nil, errors.New(fmt.Sprintf("%s arrayAdapter array is nil", adapter.field.EName))
	}
	arrV := make([]com.CalculateObj, len(adapter.field.Array))
	copy(arrV, adapter.field.Array)
	EName := adapter.field.EName
	arrType := adapter.field.ArrayType
	switch arrType {
	case conf.ReturnInt:
		var rt []int
		for k, val := range arrV {
			if len(val.UniqueId) == 0 && val.CalculateType == conf.CalculateExpr {
				val.UniqueId = uuid.GenUuidV3(val.Expr)
			}
			kName := EName + "-" + strconv.Itoa(k)
			err := adapter.fieldParse.Set(kName, adapter.params, false, &val)
			if err != nil {
				return nil, err
			}
			fv, err := adapter.fieldParse.ParseField()
			if err != nil {
				return nil, err
			}
			if fv != nil {
				toVal, ok := convert.Int(fv.GetVal())
				if ok == false {
					return nil, errors.New(fmt.Sprintf("%s arrayAdapter 值类型与arryType定义不匹配", kName))
				}
				rt = append(rt, toVal)
			} else {
				rt = append(rt, 0)
			}

		}
		return rt, nil
	case conf.ReturnFloat:
		var rt []float64
		for k, val := range arrV {
			if len(val.UniqueId) == 0 && val.CalculateType == conf.CalculateExpr {
				val.UniqueId = uuid.GenUuidV3(val.Expr)
			}
			kName := EName + "-" + strconv.Itoa(k)
			err := adapter.fieldParse.Set(kName, adapter.params, false, &val)
			if err != nil {
				return nil, err
			}
			fv, err := adapter.fieldParse.ParseField()
			if err != nil {
				return nil, err
			}
			if fv != nil {
				toVal, ok := convert.Float64(fv.GetVal())
				if ok == false {
					return nil, errors.New(fmt.Sprintf("%s arrayAdapter 值类型与arryType定义不匹配", kName))
				}
				rt = append(rt, toVal)
			} else {
				rt = append(rt, float64(0))
			}

		}
		return rt, nil
	case conf.ReturnBool:
		var rt []bool
		for k, val := range arrV {
			if len(val.UniqueId) == 0 && val.CalculateType == conf.CalculateExpr {
				val.UniqueId = uuid.GenUuidV3(val.Expr)
			}
			kName := EName + "-" + strconv.Itoa(k)
			err := adapter.fieldParse.Set(kName, adapter.params, false, &val)
			if err != nil {
				return nil, err
			}
			fv, err := adapter.fieldParse.ParseField()
			if err != nil {
				return nil, err
			}
			if fv != nil {
				toVal, ok := convert.Bool(fv.GetVal())
				if ok == false {
					return nil, errors.New(fmt.Sprintf("%s arrayAdapter 值类型与arryType定义不匹配", kName))
				}
				rt = append(rt, toVal)
			} else {
				rt = append(rt, false)
			}

		}
		return rt, nil
	case conf.ReturnStr:
		var rt []string
		for k, val := range arrV {
			if len(val.UniqueId) == 0 && val.CalculateType == conf.CalculateExpr {
				val.UniqueId = uuid.GenUuidV3(val.Expr)
			}
			kName := EName + "-" + strconv.Itoa(k)
			err := adapter.fieldParse.Set(kName, adapter.params, false, &val)
			if err != nil {
				return nil, err
			}
			fv, err := adapter.fieldParse.ParseField()
			if err != nil {
				return nil, err
			}
			if fv != nil {
				toVal, ok := convert.String(fv.GetVal())
				if ok == false {
					return nil, errors.New(fmt.Sprintf("%s arrayAdapter 值类型与arryType定义不匹配", kName))
				}
				rt = append(rt, toVal)
			} else {
				rt = append(rt, "")
			}

		}
		return rt, nil
	case conf.ReturnMap:
		var rt []map[string]interface{}
		for k, val := range arrV {
			if len(val.UniqueId) == 0 && val.CalculateType == conf.CalculateExpr {
				val.UniqueId = uuid.GenUuidV3(val.Expr)
			}
			kName := EName + "-" + strconv.Itoa(k)
			err := adapter.fieldParse.Set(kName, adapter.params, false, &val)
			if err != nil {
				return nil, err
			}
			fv, err := adapter.fieldParse.ParseField()
			if err != nil {
				return nil, err
			}
			if fv != nil {
				toVal, ok := fv.GetVal().(map[string]interface{})
				if ok == false {
					return nil, errors.New(fmt.Sprintf("%s arrayAdapter 值类型与arryType定义不匹配", kName))
				}
				rt = append(rt, toVal)
			} else {
				rt = append(rt, map[string]interface{}{})
			}

		}
		return rt, nil
	case conf.ReturnDecimal:
		var rt []decimal.Decimal
		for k, val := range arrV {
			if len(val.UniqueId) == 0 && val.CalculateType == conf.CalculateExpr {
				val.UniqueId = uuid.GenUuidV3(val.Expr)
			}
			kName := EName + "-" + strconv.Itoa(k)
			err := adapter.fieldParse.Set(kName, adapter.params, false, &val)
			if err != nil {
				return nil, err
			}
			fv, err := adapter.fieldParse.ParseField()
			if err != nil {
				return nil, err
			}
			if fv != nil {
				toVal, ok := convert.Decimal(fv.GetVal())
				if ok == false {
					return nil, errors.New(fmt.Sprintf("%s arrayAdapter 值类型与arryType定义不匹配", kName))
				}
				rt = append(rt, toVal)
			} else {
				rt = append(rt, decimal.NewFromInt32(int32(0)))
			}

		}
		return rt, nil
	case conf.ReturnDate:
		var rt []time.Time
		for k, val := range arrV {
			if len(val.UniqueId) == 0 && val.CalculateType == conf.CalculateExpr {
				val.UniqueId = uuid.GenUuidV3(val.Expr)
			}
			kName := EName + "-" + strconv.Itoa(k)
			err := adapter.fieldParse.Set(kName, adapter.params, false, &val)
			if err != nil {
				return nil, err
			}
			fv, err := adapter.fieldParse.ParseField()
			if err != nil {
				return nil, err
			}
			if fv != nil {
				toVal, ok := convert.Date(fv.GetVal())
				if ok == false {
					return nil, errors.New(fmt.Sprintf("%s arrayAdapter 值类型与arryType定义不匹配", kName))
				}
				rt = append(rt, toVal)
			} else {
				rt = append(rt, time.Time{})
			}
		}
		return rt, nil
	case conf.ReturnAny:
		var rt []interface{}
		for k, val := range arrV {
			if len(val.UniqueId) == 0 && val.CalculateType == conf.CalculateExpr {
				val.UniqueId = uuid.GenUuidV3(val.Expr)
			}
			kName := EName + "-" + strconv.Itoa(k)
			err := adapter.fieldParse.Set(kName, adapter.params, false, &val)
			if err != nil {
				return nil, err
			}
			fv, err := adapter.fieldParse.ParseField()
			if err != nil {
				return nil, err
			}
			if fv != nil {
				rt = append(rt, fv.GetVal())
			} else {
				rt = append(rt, nil)
			}
		}
		return rt, nil
	case conf.ReturnArray:
		var rt []interface{}
		for k, val := range arrV {
			if len(val.UniqueId) == 0 && val.CalculateType == conf.CalculateExpr {
				val.UniqueId = uuid.GenUuidV3(val.Expr)
			}
			kName := EName + "-" + strconv.Itoa(k)
			err := adapter.fieldParse.Set(kName, adapter.params, false, &val)
			if err != nil {
				return nil, err
			}
			fv, err := adapter.fieldParse.ParseField()
			if err != nil {
				return nil, err
			}
			if fv != nil {
				rt = append(rt, fv.GetVal())
			} else {
				rt = append(rt, nil)
			}
		}
		return rt, nil
	default:
		return nil, errors.New(fmt.Sprintf("%s arrayAdapter 不合法的数据类型", adapter.field.EName))
	}
}
