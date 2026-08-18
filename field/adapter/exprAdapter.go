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
	field      *com.Field
	params     map[string]interface{}
	expression _interface.ExprCalculate
}

func (adapter *ExprAdapter) Set(field *com.Field, params map[string]interface{}, expression _interface.ExprCalculate, f _interface.Field) {
	adapter.field = field
	adapter.params = params
	adapter.expression = expression
}

func (adapter *ExprAdapter) GetValue() (interface{}, error) {
	strExpr := adapter.field.Expr
	eName := adapter.field.EName
	uniqueId := adapter.field.UniqueId

	if len(strExpr) == 0 {
		return nil, errors.New(fmt.Sprintf("%s ExprAdapter exec err: expr is empty", adapter.field.EName))
	}
	adapter.expression.Set(strExpr, uniqueId, adapter.params, eName)
	return adapter.expression.Calculate()

}
