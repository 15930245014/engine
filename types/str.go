package types

import (
	"fmt"
	"gitlab.shudieds.com/zxh/engine/conf"
)

/**
字符串
*/

type StrExprAST struct {
	Val  string
	Name string
}

func (n StrExprAST) toStr() string {
	return fmt.Sprintf(
		"StrExprAST:%s",
		n.Name,
	)
}

func (n StrExprAST) GetVal() interface{} {
	return n.Val
}

func (n StrExprAST) GetType() string {
	return conf.ReturnStr
}
