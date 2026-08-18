package types

import "fmt"

/**
类
*/

type ClassExprAST struct {
	Val  interface{}
	Keys []string
	Name string
}

func (n ClassExprAST) toStr() string {
	return fmt.Sprintf(
		"masterExprAST:%s",
		n.Name,
	)
}

func (n ClassExprAST) GetVal() interface{} {
	return n.Val
}

func (n ClassExprAST) GetType() string {
	return "class"
}
