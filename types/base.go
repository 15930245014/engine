package types

type ExprAST interface {
	toStr() string
	GetVal() interface{}
	GetType() string
}
