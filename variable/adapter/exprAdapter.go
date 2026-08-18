package adapter

import (
	"errors"
	"fmt"
	com "gitlab.shudieds.com/mec/lib/entry/engine"
	_interface "gitlab.shudieds.com/zxh/engine/interface"
)

/**
公式adapter
*/

type ExprAdapter struct {
	variable   *com.Variable
	params     map[string]interface{}
	expression _interface.ExprCalculate
}

func (adapter *ExprAdapter) Set(variable *com.Variable, params map[string]interface{}, expression _interface.ExprCalculate, variable2 _interface.Variable) {
	adapter.variable = variable
	adapter.params = params
	adapter.expression = expression
}

func (adapter *ExprAdapter) GetValue() (interface{}, error) {
	strExpr := adapter.variable.Expr
	eName := adapter.variable.EName
	uniqueId := adapter.variable.UniqueId
	if len(strExpr) == 0 {
		return nil, errors.New(fmt.Sprintf("%s exprAdapter exec err: expr is empty", adapter.variable.EName))
	}
	adapter.expression.Set(strExpr, uniqueId, adapter.params, eName)
	return adapter.expression.Calculate()

}
