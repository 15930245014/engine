package types

import "fmt"

/**
业务对照
*/

type LookupExprAST struct {
	Val  interface{}
	Keys []string
	Name string
}

func (n LookupExprAST) toStr() string {
	return fmt.Sprintf(
		"masterExprAST:%s",
		n.Name,
	)
}

func (n LookupExprAST) GetVal() interface{} {
	return n.Val
}

func (n LookupExprAST) GetType() string {
	return "lookup"
}
