package types

import (
	"fmt"
)

/**
组合型
*/

type BinaryExprAST struct {
	Op string
	Lhs,
	Rhs ExprAST
}

func (b BinaryExprAST) toStr() string {
	return fmt.Sprintf(
		"BinaryExprAST: (%s %s %s)",
		b.Op,
		b.Lhs.toStr(),
		b.Rhs.toStr(),
	)
}

func (b BinaryExprAST) GetVal() interface{} {
	return nil
}

func (b BinaryExprAST) GetType() string {
	return ""
}
