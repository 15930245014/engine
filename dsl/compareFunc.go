package dsl

import (
	"gitlab.shudieds.com/mec/lib/utils/convert"
	"gitlab.shudieds.com/zxh/engine/types"
)

/*
*

	比较函数
*/
func (a *AST) defCompare(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}
	return convert.Compare(v1, v2)
}

/*
*

	比较函数
*/
func (a *AST) defEqual(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}
	res, ok := convert.Compare(v1, v2)
	if ok != nil {
		return nil, ok
	}

	return res == 0, nil
}
