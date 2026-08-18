package types

import (
	"fmt"
	"github.com/shopspring/decimal"
	"gitlab.shudieds.com/zxh/engine/conf"
)

/**
  decimal类型
*/

type DecimalExprAST struct {
	Val  decimal.Decimal
	Name string
}

func (n DecimalExprAST) toStr() string {
	return fmt.Sprintf(
		"DecimalExprAST:%s",
		n.Name,
	)
}

func (n DecimalExprAST) GetVal() interface{} {
	return n.Val
}
func (n DecimalExprAST) GetType() string {
	return conf.ReturnDecimal
}
