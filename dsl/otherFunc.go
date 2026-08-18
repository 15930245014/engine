package dsl

import (
	"errors"
	"github.com/dop251/goja"
	"gitlab.shudieds.com/mec/lib/utils/convert"
	"gitlab.shudieds.com/zxh/engine/types"
	"strconv"
)

/*
*

	是否为整形
*/
func (a *AST) jsExec(expr ...types.ExprAST) (interface{}, error) {
	if len(expr) == 0 {
		return nil, errors.New("jsExec 参数不能为空!")
	}
	jsExpr := expr[len(expr)-1]

	runV, err := a.ExprASTResult(jsExpr)
	if err != nil {
		return nil, errors.New("jsExec 执行错误：" + err.Error())
	}

	//转字符串
	runStr, bl := convert.String(runV)
	if bl == false {
		return nil, errors.New("jsExec 执行错误，无效的js函数")
	}

	vm := goja.New()
	_, err = vm.RunString(runStr)
	if err != nil {
		return nil, errors.New("jsExec 执行错误，无效的js函数:" + err.Error())
	}

	main, ok := goja.AssertFunction(vm.Get("main"))
	if !ok {
		return nil, errors.New("jsExec 执行错误,不存在main函数")
	}

	for i := 0; i < len(expr)-1; i++ {
		//解析
		v, err := a.ExprASTResult(expr[i])
		if err != nil {
			return nil, errors.New("jsExec 执行错误：" + err.Error())
		}
		err = vm.Set("$"+strconv.Itoa(i+1), v)
		if err != nil {
			return nil, errors.New("jsExec 执行错误：" + err.Error())
		}
	}
	res, err := main(goja.Undefined())
	if err != nil {
		return nil, errors.New("jsExec 执行错误:" + err.Error())
	}
	return res.Export(), nil
}
