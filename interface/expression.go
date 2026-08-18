package _interface

type ExprCalculate interface {
	Set(string, string, map[string]interface{}, string)
	Calculate() (interface{}, error)
}
