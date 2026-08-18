package api

import (
	"errors"
	com "gitlab.shudieds.com/mec/lib/entry/engine"
	"gitlab.shudieds.com/mec/lib/utils/convert"
	"gitlab.shudieds.com/mec/lib/utils/uuid"
	"gitlab.shudieds.com/zxh/engine/dsl"
	"gitlab.shudieds.com/zxh/engine/field"
	"gitlab.shudieds.com/zxh/engine/variable"
	"golang.org/x/net/context"
)

type ApiDslCalculate struct {
	fieldCalculate    *field.FieldCalculate
	variableCalculate *variable.VariableCalculate
	exprCalculate     *dsl.ExprCalculate
	fields            []*com.Field
	variables         []*com.Variable
	preCond           string
}

func NewApiDslCalculate(c context.Context) *ApiDslCalculate {
	apiDslCalculate := new(ApiDslCalculate)
	apiDslCalculate.fieldCalculate = field.NewFieldCalculate(c)
	apiDslCalculate.exprCalculate = dsl.NewExprCalculate(c)
	apiDslCalculate.variableCalculate = variable.NewVariableCalculate(c)

	apiDslCalculate.exprCalculate.SetFieldCalculate(apiDslCalculate.fieldCalculate)
	apiDslCalculate.exprCalculate.SetVariableCalculate(apiDslCalculate.variableCalculate)
	apiDslCalculate.fieldCalculate.SetExprCalculate(apiDslCalculate.exprCalculate)
	apiDslCalculate.variableCalculate.SetExprCalculate(apiDslCalculate.exprCalculate)

	return apiDslCalculate
}

func (a *ApiDslCalculate) Register(variables []*com.Variable, fields []*com.Field, preCond string) {
	a.variableCalculate.RegisterVariables(variables)
	a.fieldCalculate.RegisterFields(fields)
	a.fields = fields
	a.variables = variables
	a.preCond = preCond
}

func (a *ApiDslCalculate) PreCondition(params map[string]interface{}) (bool, error) {
	if len(a.preCond) == 0 {
		return true, nil
	}
	a.exprCalculate.Set(a.preCond, uuid.GenUuidV4(), params, uuid.GenUuidV4())
	val, err := a.exprCalculate.Calculate()
	if err != nil {
		return false, err
	}
	bl, _ := convert.Bool(val)
	return bl, nil
}

func (a *ApiDslCalculate) ParseFields(params map[string]interface{}) (map[string]interface{}, error) {
	result := map[string]interface{}{}
	//解析fields
	for _, f := range a.fields {
		val, err := a.ParseField(f, params)
		if err != nil {
			return result, errors.New("字段=" + f.EName + "计算失败：" + err.Error())
		} else {
			result[f.EName] = val
		}
	}
	return result, nil
}
func (a *ApiDslCalculate) ParseFieldsWithC(c context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	result := map[string]interface{}{}
	a.fieldCalculate.SetC(c)
	a.exprCalculate.SetC(c)
	a.exprCalculate.SetC(c)
	//解析fields
	for _, f := range a.fields {
		val, err := a.ParseField(f, params)
		if err != nil {
			return result, errors.New("字段=" + f.EName + "计算失败：" + err.Error())
		} else {
			result[f.EName] = val
		}
	}
	return result, nil
}
func (a *ApiDslCalculate) ParseField(f *com.Field, params map[string]interface{}) (interface{}, error) {
	result := map[string]interface{}{}
	_ = a.fieldCalculate.Set(f.EName, params, true, nil)
	expr, err := a.fieldCalculate.ParseField()
	if err != nil {
		return result, err
	}
	return expr.GetVal(), nil
}

func (a *ApiDslCalculate) ClearCache(ignoreVariables []string) {
	mp := make(map[string]bool)
	for i := 0; i < len(ignoreVariables); i++ {
		mp[ignoreVariables[i]] = true
	}
	a.variableCalculate.ClearAllCache(mp)
	a.fieldCalculate.ClearAllCache()
	a.exprCalculate.ClearCache()

}
func (a *ApiDslCalculate) GetVariableCalculate() *variable.VariableCalculate {
	return a.variableCalculate
}

func (a *ApiDslCalculate) GetFieldCalculate() *field.FieldCalculate {
	return a.fieldCalculate
}

func (a *ApiDslCalculate) GetExprCalculate() *dsl.ExprCalculate {
	return a.exprCalculate
}

func (a *ApiDslCalculate) ClearFieldCache(fName string) {
	a.fieldCalculate.ClearFieldCache(fName)
}

func (a *ApiDslCalculate) ClearVariableCache(eName string) {
	a.variableCalculate.ClearVariableCache(eName)
}

func (a *ApiDslCalculate) GetAstCache() interface{} {
	return a.exprCalculate.GetAstCache()
}
func (a *ApiDslCalculate) SetAstCache(astC interface{}) {
	a.exprCalculate.SetAstCache(astC)
}
