package types

import (
	"fmt"
	"gitlab.shudieds.com/zxh/engine/conf"
)

/**
混合型
*/

type AnyExprAST struct {
	Name string
	Val  interface{}
}

func (n AnyExprAST) toStr() string {
	return fmt.Sprintf(
		"StrExprAST:%s",
		n.Name,
	)
}
func (n AnyExprAST) GetVal() interface{} {
	return n.Val
}
func (n AnyExprAST) GetType() string {
	return conf.ReturnAny
}
