package types

import (
	"fmt"
	"github.com/shopspring/decimal"
	"gitlab.shudieds.com/zxh/engine/conf"
	"time"
)

type ArrayExprAST struct {
	IntVal     []int
	FloatVal   []float64
	BoolVal    []bool
	StrVal     []string
	MapVal     []map[string]interface{}
	ArrVal     []interface{}
	AnyVal     []interface{}
	DecimalVal []decimal.Decimal
	DateVal    []time.Time
	Name       string
	Type       string
}

func (n ArrayExprAST) toStr() string {
	return fmt.Sprintf(
		"ListExprAST:%s",
		n.Name,
	)
}

func (n ArrayExprAST) GetVal() interface{} {
	switch n.Type {
	case conf.ReturnInt:
		return n.IntVal
	case conf.ReturnFloat:
		return n.FloatVal
	case conf.ReturnStr:
		return n.StrVal
	case conf.ReturnBool:
		return n.BoolVal
	case conf.ReturnMap:
		return n.MapVal
	case conf.ReturnDate:
		return n.DateVal
	case conf.ReturnAny:
		return n.AnyVal
	case conf.ReturnArray:
		return n.ArrVal
	default:
		return nil
	}
}
func (n ArrayExprAST) GetType() string {
	return conf.ReturnArray
}
