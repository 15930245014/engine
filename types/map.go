package types

import (
	"fmt"
	"gitlab.shudieds.com/zxh/engine/conf"
)

/**
map
*/

type MapExprAST struct {
	Val  map[string]interface{}
	Name string
}

func (n MapExprAST) toStr() string {
	return fmt.Sprintf(
		"MapExprAST:%s",
		n.Name,
	)
}
func (n MapExprAST) GetVal() interface{} {
	return n.Val
}
func (n MapExprAST) GetType() string {
	return conf.ReturnMap
}
