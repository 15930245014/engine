package types

import "fmt"

type FunCallerExprAST struct {
	Name string
	Arg  []ExprAST
}

func (n FunCallerExprAST) toStr() string {
	return fmt.Sprintf(
		"FunCallerExprAST:%s",
		n.Name,
	)
}
func (n FunCallerExprAST) GetVal() interface{} {
	return nil
}

func (n FunCallerExprAST) GetType() string {
	return ""
}
