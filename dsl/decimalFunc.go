package dsl

import (
	"errors"
	"github.com/shopspring/decimal"
	"gitlab.shudieds.com/mec/lib/utils/convert"
	"gitlab.shudieds.com/zxh/engine/types"
	"strconv"
)

/*
*
 */
func (a *AST) decimalAvg(expr ...types.ExprAST) (interface{}, error) {
	if len(expr) == 0 {
		return nil, errors.New("decimalAvg参数为空")
	}
	var numbers []decimal.Decimal

	for i := 0; i < len(expr); i++ {
		v, err := a.ExprASTResult(expr[i])
		if err != nil {
			return nil, err
		}
		val, bl := v.(decimal.Decimal)
		if bl == false {
			return nil, errors.New("decimalAvg 第" + strconv.Itoa(i) + "参数不为decimal类型")
		}
		numbers = append(numbers, val)
	}

	return decimal.Avg(numbers[0], numbers[1:]...), nil
}

/*
*
 */
func (a *AST) decimalMax(expr ...types.ExprAST) (interface{}, error) {
	if len(expr) == 0 {
		return nil, errors.New("decimalMax参数为空")
	}
	var numbers []decimal.Decimal

	for i := 0; i < len(expr); i++ {
		v, err := a.ExprASTResult(expr[i])
		if err != nil {
			return nil, err
		}
		val, bl := v.(decimal.Decimal)
		if bl == false {
			return nil, errors.New("decimalMax 第" + strconv.Itoa(i) + "参数不为decimal类型")
		}
		numbers = append(numbers, val)
	}

	return decimal.Max(numbers[0], numbers[1:]...), nil
}

/*
*
 */
func (a *AST) decimalMin(expr ...types.ExprAST) (interface{}, error) {
	if len(expr) == 0 {
		return nil, errors.New("decimalMin参数为空")
	}
	var numbers []decimal.Decimal

	for i := 0; i < len(expr); i++ {
		v, err := a.ExprASTResult(expr[i])
		if err != nil {
			return nil, err
		}
		val, bl := v.(decimal.Decimal)
		if bl == false {
			return nil, errors.New("decimalMin 第" + strconv.Itoa(i) + "参数不为decimal类型")
		}
		numbers = append(numbers, val)
	}

	return decimal.Min(numbers[0], numbers[1:]...), nil
}

/*
*
 */
func (a *AST) decimalAdd(expr ...types.ExprAST) (interface{}, error) {
	if len(expr) == 0 {
		return nil, errors.New("decimalAdd参数为空")
	}
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	res, bl := v.(decimal.Decimal)
	if bl == false {
		return nil, errors.New("decimalAdd 第1个参数不为decimal类型")
	}

	for i := 1; i < len(expr); i++ {
		v1, err1 := a.ExprASTResult(expr[i])
		if err1 != nil {
			return nil, err1
		}
		val, bl := v1.(decimal.Decimal)
		if bl == false {
			return nil, errors.New("decimalAdd 第" + strconv.Itoa(i) + " 参数不为decimal类型")
		}
		res = res.Add(val)
	}

	return res, nil
}

/*
*
 */
func (a *AST) decimalSub(expr ...types.ExprAST) (interface{}, error) {
	if len(expr) == 0 {
		return nil, errors.New("decimalSub参数为空")
	}
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	res, bl := v.(decimal.Decimal)
	if bl == false {
		return nil, errors.New("decimalSub 第1个参数不为decimal类型")
	}

	for i := 1; i < len(expr); i++ {
		v1, err1 := a.ExprASTResult(expr[i])
		if err1 != nil {
			return nil, err1
		}
		val, bl := v1.(decimal.Decimal)
		if bl == false {
			return nil, errors.New("decimalSub 第" + strconv.Itoa(i) + " 参数不为decimal类型")
		}
		res = res.Sub(val)
	}

	return res, nil
}

/*
*
 */
func (a *AST) decimalDiv(expr ...types.ExprAST) (interface{}, error) {
	if len(expr) == 0 {
		return nil, errors.New("decimalDiv参数为空")
	}
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	res, bl := v.(decimal.Decimal)
	if bl == false {
		return nil, errors.New("decimalDiv 第1个参数不为decimal类型")
	}

	for i := 1; i < len(expr); i++ {
		v1, err1 := a.ExprASTResult(expr[i])
		if err1 != nil {
			return nil, err1
		}
		val, bl := v1.(decimal.Decimal)
		if bl == false {
			return nil, errors.New("decimalDiv 第" + strconv.Itoa(i) + " 参数不为decimal类型")
		}
		if val.IsZero() {
			return nil, errors.New("decimalDiv 第" + strconv.Itoa(i) + " 参数不能为0")
		}
		res = res.Div(val)
	}

	return res, nil
}

/*
*
 */
func (a *AST) decimalMul(expr ...types.ExprAST) (interface{}, error) {
	if len(expr) == 0 {
		return nil, errors.New("decimalMul参数为空")
	}
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	res, bl := v.(decimal.Decimal)
	if bl == false {
		return nil, errors.New("decimalMul 第1个参数不为decimal类型")
	}

	for i := 1; i < len(expr); i++ {
		v1, err1 := a.ExprASTResult(expr[i])
		if err1 != nil {
			return nil, err1
		}
		val, bl := v1.(decimal.Decimal)
		if bl == false {
			return nil, errors.New("decimalMul 第" + strconv.Itoa(i) + " 参数不为decimal类型")
		}
		res = res.Mul(val)
	}

	return res, nil
}

/*
*
 */
func (a *AST) decimalAbs(expr ...types.ExprAST) (interface{}, error) {
	if len(expr) == 0 {
		return nil, errors.New("decimalAbs参数为空")
	}
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	res, bl := v.(decimal.Decimal)
	if bl == false {
		return nil, errors.New("decimalAbs 参数应为decimal类型")
	}

	return res.Abs(), nil
}

/*
*
 */
func (a *AST) decimalCeil(expr ...types.ExprAST) (interface{}, error) {
	if len(expr) == 0 {
		return nil, errors.New("decimalCeil参数为空")
	}
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	res, bl := v.(decimal.Decimal)
	if bl == false {
		return nil, errors.New("decimalCeil 参数应为decimal类型")
	}

	return res.Ceil(), nil
}

/*
*
 */
func (a *AST) decimalFloor(expr ...types.ExprAST) (interface{}, error) {
	if len(expr) == 0 {
		return nil, errors.New("decimalCeil参数为空")
	}
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	res, bl := v.(decimal.Decimal)
	if bl == false {
		return nil, errors.New("decimalCeil 参数应为decimal类型")
	}

	return res.Floor(), nil
}

/*
*
 */
func (a *AST) decimalCmp(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	res, bl := v.(decimal.Decimal)
	if bl == false {
		return nil, errors.New("decimalCeil 参数1应为decimal类型")
	}
	v2, err2 := a.ExprASTResult(expr[1])
	if err2 != nil {
		return nil, err2
	}

	res2, bl := v2.(decimal.Decimal)
	if bl == false {
		return nil, errors.New("decimalCeil 参数2应为decimal类型")
	}

	return res.Cmp(res2), nil
}

/*
*
 */
func (a *AST) decimalEqual(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	res, bl := v.(decimal.Decimal)
	if bl == false {
		return nil, errors.New("decimalCeil 参数1应为decimal类型")
	}
	v2, err2 := a.ExprASTResult(expr[1])
	if err2 != nil {
		return nil, err2
	}

	res2, bl := v2.(decimal.Decimal)
	if bl == false {
		return nil, errors.New("decimalCeil 参数2应为decimal类型")
	}

	return res.Equal(res2), nil
}

/*
*
 */
func (a *AST) decimalIsZero(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	res, bl := v.(decimal.Decimal)
	if bl == false {
		return nil, errors.New("decimalIsZero 参数应为decimal类型")
	}

	return res.IsZero(), nil
}

/*
* a > b
 */
func (a *AST) decimalRt(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	res, bl := v.(decimal.Decimal)
	if bl == false {
		return nil, errors.New("decimalCeil 参数1应为decimal类型")
	}
	v2, err2 := a.ExprASTResult(expr[1])
	if err2 != nil {
		return nil, err2
	}

	res2, bl := v2.(decimal.Decimal)
	if bl == false {
		return nil, errors.New("decimalCeil 参数2应为decimal类型")
	}

	// a > b  等价于 b < a

	return res2.LessThan(res), nil
}

/*
*
 */
func (a *AST) decimalRte(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	res, bl := v.(decimal.Decimal)
	if bl == false {
		return nil, errors.New("decimalCeil 参数1应为decimal类型")
	}
	v2, err2 := a.ExprASTResult(expr[1])
	if err2 != nil {
		return nil, err2
	}

	res2, bl := v2.(decimal.Decimal)
	if bl == false {
		return nil, errors.New("decimalCeil 参数2应为decimal类型")
	}

	// a >= b 等价于  a不小于b

	isLte := res.LessThan(res2)
	if isLte == true { // a < b true ---- a >= b false
		return false, nil
	}
	return true, nil
}

/*
*
 */
func (a *AST) decimalLt(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	res, bl := v.(decimal.Decimal)
	if bl == false {
		return nil, errors.New("decimalCeil 参数1应为decimal类型")
	}
	v2, err2 := a.ExprASTResult(expr[1])
	if err2 != nil {
		return nil, err2
	}

	res2, bl := v2.(decimal.Decimal)
	if bl == false {
		return nil, errors.New("decimalCeil 参数2应为decimal类型")
	}

	return res.LessThan(res2), nil
}

/*
*
 */
func (a *AST) decimalLte(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	res, bl := v.(decimal.Decimal)
	if bl == false {
		return nil, errors.New("decimalCeil 参数1应为decimal类型")
	}
	v2, err2 := a.ExprASTResult(expr[1])
	if err2 != nil {
		return nil, err2
	}

	res2, bl := v2.(decimal.Decimal)
	if bl == false {
		return nil, errors.New("decimalCeil 参数2应为decimal类型")
	}

	return res.LessThanOrEqual(res2), nil
}

/*
*
 */
func (a *AST) decimalMod(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	res, bl := v.(decimal.Decimal)
	if bl == false {
		return nil, errors.New("decimalCeil 参数1应为decimal类型")
	}
	v2, err2 := a.ExprASTResult(expr[1])
	if err2 != nil {
		return nil, err2
	}

	res2, bl := v2.(decimal.Decimal)
	if bl == false {
		return nil, errors.New("decimalCeil 参数2应为decimal类型")
	}

	return res.Mod(res2), nil
}

/*
*
 */
func (a *AST) decimalNeg(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	res, bl := v.(decimal.Decimal)
	if bl == false {
		return nil, errors.New("decimalIsZero 参数应为decimal类型")
	}

	return res.Neg(), nil
}

/*
*
 */
func (a *AST) decimalPow(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	res, bl := v.(decimal.Decimal)
	if bl == false {
		return nil, errors.New("decimalCeil 参数1应为decimal类型")
	}
	v2, err2 := a.ExprASTResult(expr[1])
	if err2 != nil {
		return nil, err2
	}

	res2, bl := v2.(decimal.Decimal)
	if bl == false {
		return nil, errors.New("decimalCeil 参数2应为decimal类型")
	}

	return res.Pow(res2), nil
}

/*
*
 */
func (a *AST) decimalRound(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	res, bl := v.(decimal.Decimal)
	if bl == false {
		return nil, errors.New("decimalCeil 参数1应为decimal类型")
	}
	v2, err2 := a.ExprASTResult(expr[1])
	if err2 != nil {
		return nil, err2
	}

	res2, bl := convert.Int(v2)
	if bl == false {
		return nil, errors.New("decimalCeil 参数2应为整数类型")
	}

	return res.Round(int32(res2)), nil
}
