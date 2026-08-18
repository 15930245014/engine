package dsl

import (
	"errors"
	"gitlab.shudieds.com/mec/lib/pkg/sd/class/comb"
	"gitlab.shudieds.com/mec/lib/utils/convert"
	"gitlab.shudieds.com/zxh/engine/conf"
	"gitlab.shudieds.com/zxh/engine/types"
)

func (a *AST) GetEsComb(expr ...types.ExprAST) (interface{}, error) {
	//获取主数据编号
	cName := a.classParams["__C_NAME__"]
	if cName == nil {
		return nil, errors.New("GetLookupComb")
	}
	if cName != conf.CLASS_COMB {
		return nil, errors.New("GetLookupComb 无效的对象名称")
	}

	//来源系统
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, errors.New("GetLookupComb 参数1解析错误：" + err.Error())
	}
	sourceSys, ok := convert.String(v1)
	if !ok {
		return nil, errors.New("GetLookupComb  参数1返回值无法转换为str")
	}

	//获取平台
	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, errors.New("GetLookupComb 参数2解析错误：" + err.Error())
	}
	platformID, ok := convert.String(v2)
	if !ok {
		return nil, errors.New("GetLookupComb  参数2返回值无法转换为str")
	}

	//店铺
	v3, err := a.ExprASTResult(expr[2])
	if err != nil {
		return nil, errors.New("GetLookupComb 参数3解析错误：" + err.Error())
	}
	shopNo, ok := convert.String(v3)
	if !ok {
		return nil, errors.New("GetLookupComb  参数3返回值无法转换为str")
	}

	//获取msku
	v4, err := a.ExprASTResult(expr[3])
	if err != nil {
		return nil, errors.New("GetLookupComb 参数4解析错误：" + err.Error())
	}
	mSkuCode, ok := convert.String(v4)
	if !ok {
		return nil, errors.New("GetLookupComb  参数4返回值无法转换为str")
	}

	//获取sku
	v5, err := a.ExprASTResult(expr[4])
	if err != nil {
		return nil, errors.New("GetLookupComb 参数5解析错误：" + err.Error())
	}
	skuCode, ok := convert.String(v5)
	if !ok {
		return nil, errors.New("GetLookupComb  参数5返回值无法转换为str")
	}

	res, err := comb.GetEsComb(a.c, sourceSys, platformID, shopNo, mSkuCode, skuCode)
	if err != nil {
		return nil, errors.New("GetLookupComb 执行错误：" + err.Error())
	}
	return res, nil
}
