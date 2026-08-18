package dsl

import (
	"errors"
	"gitlab.shudieds.com/mec/lib/utils/convert"
	"gitlab.shudieds.com/zxh/engine/types"
	"strings"
)

/*
*
地位到高位x,x
*/
func (a *AST) bit(expr ...types.ExprAST) (interface{}, error) {
	var rtArr []string
	for i := 0; i < len(expr); i++ {
		v, err := a.ExprASTResult(expr[i])
		if err != nil {
			return nil, err
		}
		str, ok := convert.String(v)
		if !ok {
			return nil, errors.New("bitNum 参数返回值无法转换为str")
		}
		rtArr = append(rtArr, str)
	}
	return strings.Join(rtArr, ""), nil
}
