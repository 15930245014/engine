package types

import "fmt"

/**
业务主数据
*/

type MasterExprAST struct {
	Val     interface{}
	Keys []string
	Name    string
}

func (n MasterExprAST) toStr() string {
	return fmt.Sprintf(
		"masterExprAST:%s",
		n.Name,
	)
}

func (n MasterExprAST) GetVal() interface{} {
	return n.Val
}

func (n MasterExprAST) GetType() string {
	return "master"
}
