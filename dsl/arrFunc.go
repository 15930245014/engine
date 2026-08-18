package dsl

import (
	"errors"
	"github.com/shopspring/decimal"
	"gitlab.shudieds.com/mec/lib/consts"
	"gitlab.shudieds.com/mec/lib/utils/convert"
	"gitlab.shudieds.com/zxh/engine/types"
	eUtils "gitlab.shudieds.com/zxh/engine/utils"
	"time"
)

func (a *AST) defArrLen(expr ...types.ExprAST) (interface{}, error) {
	val, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	t, bl := eUtils.ArrType(val, 1)
	if bl == false {
		return nil, errors.New("arrLen 参数错误!")
	}
	toV, ok := eUtils.ToTypeVal(val, t)
	if !ok {
		return nil, errors.New("arrLen 参数错误!")
	}

	switch t {
	case consts.ARR_MAP:
		return eUtils.ArrLen(toV.([]map[string]interface{}))
	case consts.ARR_INT:
		return eUtils.ArrLen(toV.([]int))
	case consts.ARR_STR:
		return eUtils.ArrLen(toV.([]string))
	case consts.ARR_DECIMAL:
		return eUtils.ArrLen(toV.([]decimal.Decimal))
	case consts.ARR_ANY:
		return eUtils.ArrLen(toV.([]interface{}))
	case consts.ARR_BOOL:
		return eUtils.ArrLen(toV.([]bool))
	case consts.ARR_DATE:
		return eUtils.ArrLen(toV.([]time.Time))
	case consts.ARR_FLOAT:
		return eUtils.ArrLen(toV.([]float64))
	default:
		return nil, errors.New("arrLen 不合法的数组类型")
	}
}

func (a *AST) defArrIndex(expr ...types.ExprAST) (interface{}, error) {
	val, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	t, bl := eUtils.ArrType(val, 1)
	if bl == false {
		return nil, errors.New("arrIndex 参数1错误!")
	}

	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}

	i, bl2 := convert.Int(v2)
	if bl2 == false {
		return nil, errors.New("arrIndex 参数2错误!")
	}

	toV, ok := eUtils.ToTypeVal(val, t)
	if !ok {
		return nil, errors.New("arrIndex 参数错误:无法转换成数组t=" + t)
	}

	switch t {
	case consts.ARR_MAP:
		return eUtils.ArrIndex(toV.([]map[string]interface{}), i)
	case consts.ARR_INT:
		return eUtils.ArrIndex(toV.([]int), i)
	case consts.ARR_STR:
		return eUtils.ArrIndex(toV.([]string), i)
	case consts.ARR_ANY:
		return eUtils.ArrIndex(toV.([]interface{}), i)
	case consts.ARR_DECIMAL:
		return eUtils.ArrIndex(toV.([]decimal.Decimal), i)
	case consts.ARR_FLOAT:
		return eUtils.ArrIndex(toV.([]float64), i)
	case consts.ARR_BOOL:
		return eUtils.ArrIndex(toV.([]bool), i)
	case consts.ARR_DATE:
		return eUtils.ArrIndex(toV.([]time.Time), i)
	default:
		return nil, errors.New("arrIndex 不合法的数组类型t=" + t)
	}
}

func (a *AST) arrFirst(expr ...types.ExprAST) (interface{}, error) {
	val, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	t, bl := eUtils.ArrType(val, 1)
	if bl == false {
		return nil, errors.New("arrFirst 参数错误!")
	}
	toV, ok := eUtils.ToTypeVal(val, t)
	if !ok {
		return nil, errors.New("arrFirst 参数错误!无法转换成数组t=" + t)
	}
	switch t {
	case consts.ARR_MAP:
		return eUtils.ArrFirst(toV.([]map[string]interface{}))
	case consts.ARR_INT:
		return eUtils.ArrFirst(toV.([]int))
	case consts.ARR_STR:
		return eUtils.ArrFirst(toV.([]string))
	case consts.ARR_ANY:
		return eUtils.ArrFirst(toV.([]interface{}))
	case consts.ARR_DECIMAL:
		return eUtils.ArrFirst(toV.([]decimal.Decimal))
	case consts.ARR_FLOAT:
		return eUtils.ArrFirst(toV.([]float64))
	case consts.ARR_BOOL:
		return eUtils.ArrFirst(toV.([]bool))
	case consts.ARR_DATE:
		return eUtils.ArrFirst(toV.([]time.Time))
	default:
		return nil, errors.New("arrFirst 不合法的数组类型")
	}
}

func (a *AST) arrLast(expr ...types.ExprAST) (interface{}, error) {
	val, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	t, bl := eUtils.ArrType(val, 1)
	if bl == false {
		return nil, errors.New("arrLast 参数错误!")
	}
	toV, ok := eUtils.ToTypeVal(val, t)
	if !ok {
		return nil, errors.New("arrLast 参数错误!无法转换成数组t=" + t)
	}
	switch t {
	case consts.ARR_MAP:
		return eUtils.ArrLast(toV.([]map[string]interface{}))
	case consts.ARR_INT:
		return eUtils.ArrLast(toV.([]int))
	case consts.ARR_STR:
		return eUtils.ArrLast(toV.([]string))
	case consts.ARR_ANY:
		return eUtils.ArrLast(toV.([]interface{}))
	case consts.ARR_DECIMAL:
		return eUtils.ArrLast(toV.([]decimal.Decimal))
	case consts.ARR_FLOAT:
		return eUtils.ArrLast(toV.([]float64))
	case consts.ARR_BOOL:
		return eUtils.ArrLast(toV.([]bool))
	case consts.ARR_DATE:
		return eUtils.ArrLast(toV.([]time.Time))
	default:
		return nil, errors.New("arrLast 不合法的数组类型")
	}

}

func (a *AST) arrSlice(expr ...types.ExprAST) (interface{}, error) {
	val, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	t, bl := eUtils.ArrType(val, 1)
	if bl == false {
		return nil, errors.New("arrSlice 参数错误!")
	}
	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}
	v3, err := a.ExprASTResult(expr[2])
	if err != nil {
		return nil, err
	}
	s, bl2 := convert.Int(v2)
	e, bl3 := convert.Int(v3)
	if bl2 == false || bl3 == false {
		return nil, errors.New("arrSlice 参数错误!")
	}
	toV, ok := eUtils.ToTypeVal(val, t)
	if !ok {
		return nil, errors.New("arrSlice 参数错误!无法转化成数组t=" + t)
	}
	switch t {
	case consts.ARR_MAP:
		return eUtils.ArrSlice(toV.([]map[string]interface{}), s, e)
	case consts.ARR_INT:
		return eUtils.ArrSlice(toV.([]int), s, e)
	case consts.ARR_STR:
		return eUtils.ArrSlice(toV.([]string), s, e)
	case consts.ARR_ANY:
		return eUtils.ArrSlice(toV.([]any), s, e)
	case consts.ARR_DECIMAL:
		return eUtils.ArrSlice(toV.([]decimal.Decimal), s, e)
	case consts.ARR_FLOAT:
		return eUtils.ArrSlice(toV.([]float64), s, e)
	case consts.ARR_BOOL:
		return eUtils.ArrSlice(toV.([]bool), s, e)
	case consts.ARR_DATE:
		return eUtils.ArrSlice(toV.([]time.Time), s, e)
	default:
		return nil, errors.New("arrSlice 不合法的数组类型")
	}

}

func (a *AST) arrAppend(expr ...types.ExprAST) (interface{}, error) {
	val, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	v, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}
	t, bl := eUtils.ArrType(val, 2)
	if bl == false {
		return nil, errors.New("arrAppend 参数1错误!")
	}

	toV, ok := eUtils.ToTypeVal(val, t)
	if !ok {
		return nil, errors.New("arrAppend 参数错误!无法转换成数组t=" + t)
	}
	switch t {
	case consts.ARR_MAP:
		return eUtils.ArrAppend(toV.([]map[string]interface{}), v, consts.MAP)
	case consts.ARR_INT:
		return eUtils.ArrAppend(toV.([]int), v, consts.INT)
	case consts.ARR_STR:
		return eUtils.ArrAppend(toV.([]string), v, consts.STR)
	case consts.ARR_ANY:
		return eUtils.ArrAppend(toV.([]any), v, consts.ANY)
	case consts.ARR_DECIMAL:
		return eUtils.ArrAppend(toV.([]decimal.Decimal), v, consts.DECIMAL)
	case consts.ARR_FLOAT:
		return eUtils.ArrAppend(toV.([]float64), v, consts.FLOAT)
	case consts.ARR_BOOL:
		return eUtils.ArrAppend(toV.([]bool), v, consts.BOOL)
	case consts.ARR_DATE:
		return eUtils.ArrAppend(toV.([]time.Time), v, consts.DATE)
	default:
		return nil, errors.New("arrAppend 不合法的数组类型")
	}
}

func (a *AST) arrUnshift(expr ...types.ExprAST) (interface{}, error) {
	val, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	v, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}

	t, bl := eUtils.ArrType(val, 2)
	if bl == false {
		return nil, errors.New("arrUnshift 参数错误!")
	}

	toV, ok := eUtils.ToTypeVal(val, t)
	if !ok {
		return nil, errors.New("arrUnshift 参数错误!无法转换成数组t=" + t)
	}
	switch t {
	case consts.ARR_MAP:
		return eUtils.ArrUnshift(toV.([]map[string]interface{}), v, consts.MAP)
	case consts.ARR_INT:
		return eUtils.ArrUnshift(toV.([]int), v, consts.INT)
	case consts.ARR_STR:
		return eUtils.ArrUnshift(toV.([]string), v, consts.STR)
	case consts.ARR_ANY:
		return eUtils.ArrUnshift(toV.([]any), v, consts.ANY)
	case consts.ARR_DECIMAL:
		return eUtils.ArrUnshift(toV.([]decimal.Decimal), v, consts.DECIMAL)
	case consts.ARR_FLOAT:
		return eUtils.ArrUnshift(toV.([]float64), v, consts.FLOAT)
	case consts.ARR_BOOL:
		return eUtils.ArrUnshift(toV.([]bool), v, consts.BOOL)
	case consts.ARR_DATE:
		return eUtils.ArrUnshift(toV.([]time.Time), v, consts.DATE)
	default:
		return nil, errors.New("arrUnshift 不合法的数组类型")
	}
}

func (a *AST) arrMerge(expr ...types.ExprAST) (interface{}, error) {
	val1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	val2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}
	t1, bl1 := eUtils.ArrType(val1, 2)
	t2, bl2 := eUtils.ArrType(val2, 2)
	if !bl1 || !bl2 {
		return nil, errors.New("arrMerge 参数错误!")
	}
	if t1 != t2 {
		return nil, errors.New("arrMerge 参数错误：数组类型不一致")
	}
	toV1, ok := eUtils.ToTypeVal(val1, t1)
	if !ok {
		return nil, errors.New("arrMerge 参数错误!无法转换成数组t=" + t1)
	}
	toV2, ok := eUtils.ToTypeVal(val2, t2)
	if !ok {
		return nil, errors.New("arrMerge 参数错误!无法转换成数组t=" + t2)
	}
	switch t1 {
	case consts.ARR_MAP:
		return eUtils.ArrMerge(toV1.([]map[string]interface{}), toV2.([]map[string]interface{}))
	case consts.ARR_INT:
		return eUtils.ArrMerge(toV1.([]int), toV2.([]int))
	case consts.ARR_STR:
		return eUtils.ArrMerge(toV1.([]string), toV2.([]string))
	case consts.ARR_ANY:
		return eUtils.ArrMerge(toV1.([]any), toV2.([]any))
	case consts.ARR_DECIMAL:
		return eUtils.ArrMerge(toV1.([]decimal.Decimal), toV2.([]decimal.Decimal))
	case consts.ARR_FLOAT:
		return eUtils.ArrMerge(toV1.([]float64), toV2.([]float64))
	case consts.ARR_BOOL:
		return eUtils.ArrMerge(toV1.([]bool), toV2.([]bool))
	case consts.ARR_DATE:
		return eUtils.ArrMerge(toV1.([]time.Time), toV2.([]time.Time))
	default:
		return nil, errors.New("arrMerge 不合法的数组类型")
	}
}

func (a *AST) arrSearch(expr ...types.ExprAST) (interface{}, error) {
	val, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	t, bl := eUtils.ArrType(val, 1)
	if bl == false {
		return nil, errors.New("arrSearch 参数错误!")
	}
	toV, ok := eUtils.ToTypeVal(val, t)
	if !ok {
		return nil, errors.New("arrSearch 参数错误!无法转化成数组t=" + t)
	}

	v, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}
	switch t {
	case consts.ARR_STR:
		return eUtils.ArrSearch(toV.([]string), v, t)
	case consts.ARR_INT:
		return eUtils.ArrSearch(toV.([]int), v, t)
	case consts.ARR_DECIMAL:
		return eUtils.ArrSearch(toV.([]decimal.Decimal), v, t)
	case consts.ARR_DATE:
		return eUtils.ArrSearch(toV.([]time.Time), v, t)
	case consts.ARR_FLOAT:
		return eUtils.ArrSearch(toV.([]float64), v, t)
	case consts.ARR_BOOL:
		return eUtils.ArrSearch(toV.([]bool), v, t)
	default:
		return nil, errors.New("ArrSearch 不合法的数组类型")
	}
}

func (a *AST) arrUnique(expr ...types.ExprAST) (interface{}, error) {
	val, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	t, bl := eUtils.ArrType(val, 1)
	if bl == false {
		return nil, errors.New("arrUnique 参数错误!")
	}
	toV, ok := eUtils.ToTypeVal(val, t)
	if !ok {
		return nil, errors.New("arrUnique 参数错误!无法转换成数组t=" + t)
	}
	switch t {
	case consts.ARR_STR:
		return eUtils.ArrUnique(toV.([]string))
	case consts.ARR_INT:
		return eUtils.ArrUnique(toV.([]int))
	case consts.ARR_DATE:
		return eUtils.ArrUnique(toV.([]time.Time))
	case consts.ARR_DECIMAL:
		return eUtils.ArrUnique(toV.([]decimal.Decimal))
	case consts.ARR_FLOAT:
		return eUtils.ArrUnique(toV.([]float64))
	case consts.ARR_BOOL:
		return eUtils.ArrUnique(toV.([]bool))
	default:
		return nil, errors.New("ArrUnique 不合法的数组类型")
	}
}

func (a *AST) arrReverse(expr ...types.ExprAST) (interface{}, error) {
	val, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	t, bl := eUtils.ArrType(val, 1)
	if bl == false {
		return nil, errors.New("arrReverse 参数错误!")
	}
	toV, ok := eUtils.ToTypeVal(val, t)
	if !ok {
		return nil, errors.New("arrSearch 参数错误!无法转化成数组t=" + t)
	}
	switch t {
	case consts.ARR_MAP:
		return eUtils.ArrReverse(toV.([]map[string]interface{}))
	case consts.ARR_INT:
		return eUtils.ArrReverse(toV.([]int))
	case consts.ARR_STR:
		return eUtils.ArrReverse(toV.([]string))
	case consts.ARR_ANY:
		return eUtils.ArrReverse(toV.([]interface{}))
	case consts.ARR_DECIMAL:
		return eUtils.ArrReverse(toV.([]decimal.Decimal))
	case consts.ARR_FLOAT:
		return eUtils.ArrReverse(toV.([]float64))
	case consts.ARR_BOOL:
		return eUtils.ArrReverse(toV.([]bool))
	case consts.ARR_DATE:
		return eUtils.ArrReverse(toV.([]time.Time))
	default:
		return nil, errors.New("arrReverse 不合法的数组类型")
	}

}

func (a *AST) arrSortAsc(expr ...types.ExprAST) (interface{}, error) {
	val, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	t, bl := eUtils.ArrType(val, 1)
	if bl == false {
		return nil, errors.New("arrSortAsc 参数错误!")
	}

	toV, ok := eUtils.ToTypeVal(val, t)
	if !ok {
		return nil, errors.New("arrSortAsc 参数错误!无法转换成数组t=" + t)
	}
	switch t {
	case consts.ARR_INT:
		return eUtils.ArrSortAsc(toV.([]int))
	case consts.ARR_STR:
		return eUtils.ArrSortAsc(toV.([]string))
	case consts.ARR_DATE:
		return eUtils.ArrSortAsc(toV.([]time.Time))
	case consts.ARR_DECIMAL:
		return eUtils.ArrSortAsc(toV.([]decimal.Decimal))
	case consts.ARR_FLOAT:
		return eUtils.ArrSortAsc(toV.([]float64))
	case consts.ARR_BOOL:
		return eUtils.ArrSortAsc(toV.([]bool))
	default:
		return nil, errors.New("ArrSortAsc 不合法的数组类型")
	}
}

func (a *AST) arrSortDesc(expr ...types.ExprAST) (interface{}, error) {
	val, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	t, bl := eUtils.ArrType(val, 1)
	if bl == false {
		return nil, errors.New("arrSortAsc 参数错误!")
	}

	toV, ok := eUtils.ToTypeVal(val, t)
	if !ok {
		return nil, errors.New("arrSortAsc 参数错误!无法转化成数组t=" + t)
	}
	switch t {
	case consts.ARR_STR:
		return eUtils.ArrSortDesc(toV.([]string))
	case consts.ARR_INT:
		return eUtils.ArrSortDesc(toV.([]int))
	case consts.ARR_DATE:
		return eUtils.ArrSortDesc(toV.([]time.Time))
	case consts.ARR_DECIMAL:
		return eUtils.ArrSortDesc(toV.([]decimal.Decimal))
	case consts.ARR_FLOAT:
		return eUtils.ArrSortDesc(toV.([]float64))
	case consts.ARR_BOOL:
		return eUtils.ArrSortDesc(toV.([]bool))
	default:
		return nil, errors.New("ArrSortDesc 不合法的数组类型")
	}
}
