package _interface

import (
	com "gitlab.shudieds.com/mec/lib/entry/engine"
	"gitlab.shudieds.com/zxh/engine/types"
)

type Variable interface {
	Set(string, map[string]interface{}, bool, *com.CalculateObj) error
	ParseVariable() (types.ExprAST, error)
}
