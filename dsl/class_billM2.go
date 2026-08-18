package dsl

import (
	"errors"
	"gitlab.shudieds.com/mec/lib/pkg/sd/class/bill"
	"gitlab.shudieds.com/mec/lib/utils/convert"
	"gitlab.shudieds.com/zxh/engine/conf"
	"gitlab.shudieds.com/zxh/engine/types"
)

func (a *AST) GetEsBillM2(expr ...types.ExprAST) (interface{}, error) {
	//获取主数据编号
	cName := a.classParams["__C_NAME__"]
	if cName == nil {
		return nil, errors.New("GetEsBillM2")
	}
	if cName != conf.CLASS_BILL {
		return nil, errors.New("GetEsBillM2 无效的对象名称")
	}

	//获取店铺
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, errors.New("GetEsBillM2 参数1解析错误：" + err.Error())
	}
	shopNo, ok := convert.String(v1)
	if !ok {
		return nil, errors.New("GetEsBillM2  参数1返回值无法转换为str")
	}

	//获取账期
	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, errors.New("GetEsBillM2 参数2解析错误：" + err.Error())
	}
	period, ok := convert.String(v2)
	if !ok {
		return nil, errors.New("GetEsBillM2  参数2返回值无法转换为str")
	}

	//获取收支类型
	v3, err := a.ExprASTResult(expr[2])
	if err != nil {
		return nil, errors.New("GetEsBillM2 参数3解析错误：" + err.Error())
	}
	inoutType, ok := convert.String(v3)
	if !ok {
		return nil, errors.New("GetEsBillM2  参数3返回值无法转换为str")
	}

	//获取费用项
	v4, err := a.ExprASTResult(expr[3])
	if err != nil {
		return nil, errors.New("GetEsBillM2 参数4解析错误：" + err.Error())
	}
	item, ok := convert.String(v4)
	if !ok {
		return nil, errors.New("GetEsBillM2  参数4返回值无法转换为str")
	}

	res, err := bill.GetEsBillM2(a.c, shopNo, period, inoutType, item)
	if err != nil {
		return nil, errors.New("GetEsBillM2 执行错误：" + err.Error())
	}
	return res, nil
}
