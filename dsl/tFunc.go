package dsl

import (
	"errors"
	"fmt"
	"gitlab.shudieds.com/mec/lib/utils/convert"
	"gitlab.shudieds.com/zxh/engine/types"
	eUtils "gitlab.shudieds.com/zxh/engine/utils"
)

/*
*

	是否为整形
*/
func (a *AST) toInt(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err == nil {
		val, bl := convert.Int(v)
		if bl {
			return val, nil
		}
	}
	if len(expr) == 1 {
		return nil, errors.New("toInt转化失败")
	}
	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}
	val2, bl := convert.Int(v2)
	if bl == false {
		return nil, errors.New("toInt转化失败")
	}
	return val2, nil
}

/*
*

	转化float
*/
func (a *AST) toFloat(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err == nil {
		val, bl := convert.Float64(v)
		if bl {
			return val, nil
		}
	}

	//报错
	if len(expr) == 1 {
		return nil, errors.New("toFloat转化失败")
	}

	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}
	val2, bl := convert.Float64(v2)
	if bl == false {
		return nil, errors.New("toFloat转化失败")
	}

	return val2, nil
}

/*
*

	转化Str
*/
func (a *AST) toStr(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err == nil {
		val, bl := convert.String(v)
		if bl {
			return val, nil
		}
	}
	if len(expr) == 1 {
		return nil, fmt.Errorf("toStr转化失败")
	}

	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}
	val2, bl := convert.String(v2)
	if bl == false {
		return nil, fmt.Errorf("toStr转化失败")
	}
	return val2, nil
}

/*
*

	toBool
*/
func (a *AST) toBool(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err == nil {
		val, bl := convert.Bool(v)
		if bl {
			return val, nil
		}
	}

	if len(expr) == 1 {
		return nil, fmt.Errorf("toBool转化失败%v", v)
	}

	//兜底
	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}

	val2, bl := convert.Bool(v2)

	if bl == false {
		return nil, fmt.Errorf("toBool转化失败%v", v2)
	}

	return val2, nil
}

/*
*

	转化date
*/
func (a *AST) toDate(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err == nil {
		val, bl := convert.Date(v)
		if bl {
			return val, nil
		}
	}

	if len(expr) == 1 {
		return nil, fmt.Errorf("toDate转化失败%v", v)
	}

	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}
	val2, bl := convert.Date(v2)
	if !bl {
		return nil, fmt.Errorf("toDate转化失败%v", v2)
	}

	return val2, nil
}

/*
*

	转化为decimal
*/
func (a *AST) toDecimal(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err == nil {
		val, bl := convert.Decimal(v)
		if bl {
			return val, nil
		}
	}

	if len(expr) == 1 {
		return nil, fmt.Errorf("toDecimal转化失败:{%v} ", v)
	}

	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return false, err
	}
	val2, bl := convert.Decimal(v2)
	if bl == false {
		return nil, fmt.Errorf("toDecimal转化失败:{%v} ", v2)
	}
	return val2, nil
}

/*
*

	是否为整形
*/
func (a *AST) isInt(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return false, err
	}
	switch v.(type) {
	case int:
		return true, nil
	case int8:
		return true, nil
	case int16:
		return true, nil
	case int32:
		return true, nil
	case int64:
		return true, nil
	case uint:
		return true, nil
	case uint8:
		return true, nil
	case uint16:
		return true, nil
	case uint32:
		return true, nil
	case uint64:
		return true, nil
	}

	return false, nil

}

/*
*

	是否为float
*/
func (a *AST) isFloat(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return false, err
	}
	_, ok := convert.Float64(v)
	return ok, nil
}

/*
*

	是否为Str
*/
func (a *AST) isStr(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return false, err
	}
	_, bl := convert.String(v)
	return bl, nil
}

/*
*

	isBool
*/
func (a *AST) isBool(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return false, err
	}
	_, bl := convert.Bool(v)
	return bl, nil
}

/*
*

	是否为date
*/
func (a *AST) isDate(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return false, err
	}
	_, bl := convert.Date(v)
	return bl, nil
}

/*
*

	是否为decimal
*/
func (a *AST) isDecimal(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return false, nil
	}
	_, bl := convert.Decimal(v)
	return bl, nil
}

/*
*

	是否为
*/
func (a *AST) isObj(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return false, nil
	}
	_, bl := eUtils.ToMap(v)
	return bl, nil
}

/*
*

	是否为array
*/
func (a *AST) isArray(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return false, nil
	}
	_, bl := eUtils.ArrType(v, 1)
	return bl, nil
}

/*
*

	是否为zero
*/
func (a *AST) isZero(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return false, errors.New("isZero 参数计算失败：" + err.Error())
	}
	t, _ := convert.Compare(v, "__DSL-ZERO__")
	return t == 0, nil
}
