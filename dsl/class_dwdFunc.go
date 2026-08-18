package dsl

import (
	"errors"
	"gitlab.shudieds.com/mec/lib/consts"
	"gitlab.shudieds.com/mec/lib/pkg/sd/class/dwd"
	"gitlab.shudieds.com/mec/lib/utils/convert"
	"gitlab.shudieds.com/zxh/engine/conf"
	"gitlab.shudieds.com/zxh/engine/types"
	eUtils "gitlab.shudieds.com/zxh/engine/utils"
)

func (a *AST) GetEsDwdInformation(expr ...types.ExprAST) (interface{}, error) {
	//获取主数据编号
	cName := a.classParams["__C_NAME__"]
	if cName == nil {
		return nil, errors.New("GetEsDwdInformation主数据对象名称未指定")
	}
	if cName != conf.CLASS_DWD {
		return nil, errors.New("GetEsDwdInformation 无效的对象名称")
	}

	//获取索引名称
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, errors.New("GetEsDwdInformation 参数1解析错误：" + err.Error())
	}
	indexName, ok := convert.String(v1)
	if !ok {
		return nil, errors.New("GetEsDwdInformation  参数1返回值无法转换为str")
	}

	//获取来源系统
	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, errors.New("GetEsDwdInformation 参数2解析错误：" + err.Error())
	}
	source, ok := convert.String(v2)
	if !ok {
		return nil, errors.New("GetEsDwdInformation  参数2返回值无法转换为str")
	}

	//获取交易类型
	v3, err := a.ExprASTResult(expr[2])
	if err != nil {
		return nil, errors.New("GetEsDwdInformation 参数3解析错误：" + err.Error())
	}
	transactionType, ok := convert.String(v3)
	if !ok {
		return nil, errors.New("GetEsDwdInformation  参数3返回值无法转换为str")
	}

	//获取交易单号
	v4, err := a.ExprASTResult(expr[3])
	if err != nil {
		return nil, errors.New("GetEsDwdInformation 参数4解析错误：" + err.Error())
	}
	bizNo, ok := eUtils.ToTypeVal(v4, consts.ARR_STR)
	if !ok {
		return nil, errors.New("GetEsDwdInformation  参数4返回值无法转换为str数组")
	}

	//获取索引
	v5, err := a.ExprASTResult(expr[4])
	if err != nil {
		return nil, errors.New("GetEsDwdInformation 参数5解析错误：" + err.Error())
	}
	keys, _ := convert.String(v5)

	res, err := dwd.GetEsDwdInformation(a.c, indexName, source, transactionType, bizNo.([]string), keys)
	if err != nil {
		return nil, errors.New("GetEsDwdInformation 执行错误：" + err.Error())
	}
	return res, nil
}

func (a *AST) GetEsDwdChildInformation(expr ...types.ExprAST) (interface{}, error) {
	//获取主数据编号
	cName := a.classParams["__C_NAME__"]
	if cName == nil {
		return nil, errors.New("GetEsDwdInformation主数据对象名称未指定")
	}
	if cName != conf.CLASS_DWD {
		return nil, errors.New("GetEsDwdInformation 无效的对象名称")
	}

	//获取索引名称
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, errors.New("GetEsDwdInformation 参数1解析错误：" + err.Error())
	}
	indexName, ok := convert.String(v1)
	if !ok {
		return nil, errors.New("GetEsDwdInformation  参数1返回值无法转换为str")
	}

	//获取来源系统
	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, errors.New("GetEsDwdInformation 参数2解析错误：" + err.Error())
	}
	source, ok := convert.String(v2)
	if !ok {
		return nil, errors.New("GetEsDwdInformation  参数2返回值无法转换为str")
	}

	//获取交易类型
	v3, err := a.ExprASTResult(expr[2])
	if err != nil {
		return nil, errors.New("GetEsDwdInformation 参数3解析错误：" + err.Error())
	}
	transactionType, ok := convert.String(v3)
	if !ok {
		return nil, errors.New("GetEsDwdInformation  参数3返回值无法转换为str")
	}

	//获取交易单号
	v4, err := a.ExprASTResult(expr[3])
	if err != nil {
		return nil, errors.New("GetEsDwdInformation 参数4解析错误：" + err.Error())
	}
	childNo, ok := eUtils.ToTypeVal(v4, consts.ARR_STR)
	if !ok {
		return nil, errors.New("GetEsDwdInformation  参数4返回值无法转换为str数组")
	}

	//获取索引
	v5, err := a.ExprASTResult(expr[4])
	if err != nil {
		return nil, errors.New("GetEsDwdInformation 参数5解析错误：" + err.Error())
	}
	keys, _ := convert.String(v5)

	res, err := dwd.GetEsDwdChildInformation(a.c, indexName, source, transactionType, childNo.([]string), keys)
	if err != nil {
		return nil, errors.New("GetEsDwdInformation 执行错误：" + err.Error())
	}
	return res, nil
}
