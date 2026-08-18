package _interface

import (
	com "gitlab.shudieds.com/mec/lib/entry/engine"
	"gitlab.shudieds.com/zxh/engine/types"
)

type Field interface {
	Set(string, map[string]interface{}, bool, *com.CalculateObj) error
	ParseField() (types.ExprAST, error)
}
