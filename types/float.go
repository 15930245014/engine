package types

import (
	"fmt"
	"gitlab.shudieds.com/zxh/engine/conf"
)

/**
数字型
*/

type FloatExprAST struct {
	Val  float64
	Name string
}

func (n FloatExprAST) toStr() string {
	return fmt.Sprintf(
		"NumberExprAST:%s",
		n.Name,
	)
}
func (n FloatExprAST) GetVal() interface{} {
	return n.Val
}
func (n FloatExprAST) GetType() string {
	return conf.ReturnFloat
}
