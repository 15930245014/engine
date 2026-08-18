package types

import (
	"fmt"
	"gitlab.shudieds.com/zxh/engine/conf"
	"time"
)

/**
时间格式
*/

type DateExprAST struct {
	Val  time.Time
	Name string
}

func (n DateExprAST) toStr() string {
	return fmt.Sprintf(
		"DateExprAST:%s",
		n.Name,
	)
}

func (n DateExprAST) GetVal() interface{} {
	return n.Val
}
func (n DateExprAST) GetType() string {
	return conf.ReturnDate
}
