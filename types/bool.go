package types

import (
	"fmt"
	"gitlab.shudieds.com/zxh/engine/conf"
)

/**
bool
*/

type BoolExprAST struct {
	Val  bool
	Name string
}

func (n BoolExprAST) toStr() string {
	return fmt.Sprintf(
		"BoolExprAST:%s",
		n.Name,
	)
}
func (n BoolExprAST) GetVal() interface{} {
	return n.Val
}

func (n BoolExprAST) GetType() string {
	return conf.ReturnBool
}
