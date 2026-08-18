package dsl

import (
	"errors"
	"gitlab.shudieds.com/mec/lib/consts"
	"gitlab.shudieds.com/mec/lib/pkg/sd/class/dws"
	"gitlab.shudieds.com/zxh/engine/conf"
	"gitlab.shudieds.com/zxh/engine/types"
	eUtils "gitlab.shudieds.com/zxh/engine/utils"
)

func (a *AST) GetEsDwsStatement(expr ...types.ExprAST) (interface{}, error) {
	//获取主数据编号
	cName := a.classParams["__C_NAME__"]
	if cName == nil {
		return nil, errors.New("GetEsDwsStatement")
	}
	if cName != conf.CLASS_DWS {
		return nil, errors.New("GetEsDwsStatement 无效的对象名称")
	}

	//获取对账单UUID
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, errors.New("GetEsDwsStatement 参数1解析错误：" + err.Error())
	}
	ordSetUuids, ok := eUtils.ToTypeVal(v1, consts.ARR_STR)
	if !ok {
		return nil, errors.New("GetEsDwsStatement  参数1返回值无法转换为str数组")
	}
	res, err := dws.GetEsDwsStatement(a.c, ordSetUuids.([]string))
	if err != nil {
		return nil, errors.New("GetEsDwsStatement 执行错误：" + err.Error())
	}
	return res, nil
}
