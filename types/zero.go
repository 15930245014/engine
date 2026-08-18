package types

import (
	"fmt"
	"gitlab.shudieds.com/zxh/engine/conf"
)

/**
默认值
*/

type ZeroExprAST struct {
	Val  interface{}
	Name string
}

func (n ZeroExprAST) toStr() string {
	return fmt.Sprintf(
		"StrExprAST:%s",
		n.Val,
	)
}
func (n ZeroExprAST) GetVal() interface{} {
	return n.Val
}
func (n ZeroExprAST) GetType() string {
	return conf.ReturnZero
}
