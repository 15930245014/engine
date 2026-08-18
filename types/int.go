package types

import (
	"fmt"
	"gitlab.shudieds.com/zxh/engine/conf"
)

/**
整形
*/

type IntExprAST struct {
	Val  int
	Name string
}

func (n IntExprAST) toStr() string {
	return fmt.Sprintf(
		"NumberExprAST:%s",
		n.Name,
	)
}
func (n IntExprAST) GetVal() interface{} {
	return n.Val
}
func (n IntExprAST) GetType() string {
	return conf.ReturnInt
}
