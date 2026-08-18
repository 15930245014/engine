package types

import (
	"fmt"
	"gitlab.shudieds.com/zxh/engine/conf"
)

/**
self
*/

type SelfExprAST struct {
	Val  interface{}
	Name string
}

func (n SelfExprAST) toStr() string {
	return fmt.Sprintf(
		"SelfExprAST:%s",
		n.Val,
	)
}
func (n SelfExprAST) GetVal() interface{} {
	return n.Val
}
func (n SelfExprAST) GetType() string {
	return conf.ReturnSelf
}
