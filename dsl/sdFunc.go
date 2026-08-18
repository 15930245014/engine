package dsl

import (
	"errors"
	"gitlab.shudieds.com/mec/lib/entry/engine"
	"gitlab.shudieds.com/mec/lib/pkg/sd"
	"gitlab.shudieds.com/mec/lib/utils/convert"
	"gitlab.shudieds.com/zxh/engine/types"
)

func (a *AST) MGet(expr ...types.ExprAST) (interface{}, error) {
	//获取主数据编号
	mName := a.mParams["__M_NAME__"]
	if mName == nil {
		return nil, errors.New("MGet主数据对象名称未指定")
	}
	mNameStr, ok := convert.String(mName)
	if !ok {
		return nil, errors.New("MGet主数据对象名称不合法，应为字符串类型")
	}

	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, errors.New("MGet主数据对象名称不合法，参数解析错误：" + err.Error())
	}
	code, ok := convert.String(v)
	if !ok {
		return nil, errors.New("MGet主数据编码 参数返回值无法转换为str")
	}
	return sd.MGet(a.c, mNameStr, code)

}
func (a *AST) LGet(expr ...types.ExprAST) (interface{}, error) {
	//获取主数据编号
	mName := a.lParams["__L_NAME__"]
	if mName == nil {
		return nil, errors.New("LGet主数据对象名称未指定")
	}
	mNameStr, ok := convert.String(mName)
	if !ok {
		return nil, errors.New("LGet主数据对象名称不合法，应为字符串类型")
	}
	v0, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, errors.New("LGet 来源系统字段获取错误：" + err.Error())
	}
	sSys, ok := convert.String(v0)
	if !ok {
		return nil, errors.New("LGet来源系统返回值无法转换为str")
	}

	v1, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, errors.New("LGet 来源系统编码字段获取错误：" + err.Error())
	}
	sCode, ok := convert.String(v1)
	if !ok {
		return nil, errors.New("LGet来源系统编码返回值无法转换为str")
	}

	v2, err := a.ExprASTResult(expr[2])
	if err != nil {
		return nil, errors.New("LGet 目标系统字段获取错误：" + err.Error())
	}
	tSys, ok := convert.String(v2)
	if !ok {
		return nil, errors.New("LGet 目标系统返回值无法转换为str")
	}
	return sd.LGet(a.c, mNameStr, sSys, sCode, tSys)

}

func (a *AST) LGetV2(expr ...types.ExprAST) (interface{}, error) {
	//获取主数据编号
	mName := a.lParams["__L_NAME__"]
	if mName == nil {
		return nil, errors.New("LGetV2主数据对象名称未指定")
	}
	mNameStr, ok := convert.String(mName)
	if !ok {
		return nil, errors.New("LGetV2主数据对象名称不合法，应为字符串类型")
	}
	v0, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, errors.New("LGetV2 来源系统字段获取错误：" + err.Error())
	}
	sSys, ok := convert.String(v0)
	if !ok {
		return nil, errors.New("LGetV2 来源系统返回值无法转换为str")
	}

	v1, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, errors.New("LGetV2 来源系统编码字段获取错误：" + err.Error())
	}
	sCode, ok := convert.String(v1)
	if !ok {
		return nil, errors.New("LGetV2 来源系统编码返回值无法转换为str")
	}

	v2, err := a.ExprASTResult(expr[2])
	if err != nil {
		return nil, errors.New("LGetV2 租户字段获取错误：" + err.Error())
	}
	tent, ok := convert.String(v2)
	if !ok {
		return nil, errors.New("LGetV2 租户返回值无法转换为str")
	}

	v3, err := a.ExprASTResult(expr[3])
	if err != nil {
		return nil, errors.New("LGetV2 目标系统字段获取错误：" + err.Error())
	}
	tSys, ok := convert.String(v3)
	if !ok {
		return nil, errors.New("LGetV2 目标系统返回值无法转换为str")
	}
	return sd.LGetV2(a.c, mNameStr, sSys, sCode, tent, tSys)

}

/**
P函数
*/

func (a *AST) P(expr ...types.ExprAST) (interface{}, error) {
	v0, err := a.ExprASTResult(expr[0])
	if err != nil || v0 == nil {
		return nil, errors.New("P 函数获取错误：" + err.Error())
	}

	P, ok := v0.(*engine.TreeParams)

	if !ok {
		return nil, errors.New("P 函数 类型断言错误")
	}
	return P.P, nil
}
