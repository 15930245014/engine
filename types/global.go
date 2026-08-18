package types

import (
	"fmt"
)

/**
全局
*/

type GlobalExprAST struct {
	Val  interface{}
	Keys []string
	Name string
}

func (n GlobalExprAST) toStr() string {
	return fmt.Sprintf(
		"variableExprAST:%s",
		n.Name,
	)
}

func (n GlobalExprAST) GetVal() interface{} {
	return n.Val
}

func (n GlobalExprAST) GetType() string {
	return "global"
}
