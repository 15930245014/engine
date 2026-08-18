package types

import (
	"fmt"
)

/**
$变量
*/

type VariableExprAST struct {
	Val  interface{}
	Name string
	Keys []string
}

func (n VariableExprAST) toStr() string {
	return fmt.Sprintf(
		"variableExprAST:%s",
		n.Name,
	)
}

func (n VariableExprAST) GetVal() interface{} {
	return n.Val
}

func (n VariableExprAST) GetType() string {
	return "variable"
}
