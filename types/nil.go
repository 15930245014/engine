package types

import (
	"fmt"
	"gitlab.shudieds.com/zxh/engine/conf"
)

/**
空值
*/

type NilExprAST struct {
	Val  interface{}
	Name string
}

func (n NilExprAST) toStr() string {
	return fmt.Sprintf(
		"nilExprAST:%s",
		n.Name,
	)
}
func (n NilExprAST) GetVal() interface{} {
	return n.Val
}
func (n NilExprAST) GetType() string {
	return conf.ReturnNil
}
