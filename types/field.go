package types

import (
	"fmt"
)

/**
@field
*/

type FieldExprAST struct {
	Val  interface{}
	Name string
	Keys []string
}

func (n FieldExprAST) toStr() string {
	return fmt.Sprintf(
		"variableExprAST:%s",
		n.Name,
	)
}

func (n FieldExprAST) GetVal() interface{} {
	return n.Val
}

func (n FieldExprAST) GetType() string {
	return "variable"
}
