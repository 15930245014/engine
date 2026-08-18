package dsl

import (
	"context"
	"errors"
	"fmt"
	"gitlab.shudieds.com/zxh/engine/conf"
	_interface "gitlab.shudieds.com/zxh/engine/interface"
	"gitlab.shudieds.com/zxh/engine/types"
	"strconv"
)

/**
ast 结构体
*/

type AST struct {
	c                 context.Context
	Tokens            []Token
	source            string
	currTok           Token
	currIndex         int
	depth             int
	Err               error
	globalParams      map[string]interface{} //全局参数
	localParams       map[string]interface{} //局部参数
	mParams           map[string]interface{}
	lParams           map[string]interface{}
	classParams       map[string]interface{}
	selfParams        map[string]interface{}
	variableCalculate _interface.Variable
	fieldCalculate    _interface.Field
	defFunc           map[string]defS
	ignoreFunc        map[string]bool
	dataType          map[string]bool
	root              types.ExprAST
	exprCal           *ExprCalculate
}

/**
  生成语法树
*/

func NewAST(c context.Context, toks []Token, s string, globalParams map[string]interface{}, variableCalculate _interface.Variable, field _interface.Field) *AST {
	a := &AST{
		c:                 c,
		Tokens:            toks,
		source:            s,
		globalParams:      globalParams,
		variableCalculate: variableCalculate,
		fieldCalculate:    field,
		localParams:       map[string]interface{}{},
		mParams:           map[string]interface{}{},
		lParams:           map[string]interface{}{},
		selfParams:        map[string]interface{}{},
		classParams:       map[string]interface{}{},
	}
	if a.Tokens == nil || len(a.Tokens) == 0 {
		a.Err = errors.New("empty token")
	} else {
		a.currIndex = 0
		a.currTok = a.Tokens[0]
	}
	a.initFunc()
	a.dataType = map[string]bool{
		conf.ReturnInt:     true,
		conf.ReturnStr:     true,
		conf.ReturnDate:    true,
		conf.ReturnFloat:   true,
		conf.ReturnDecimal: true,
		conf.ReturnBool:    true,
		conf.ReturnMap:     true,
		conf.ReturnArray:   true,
		conf.ReturnAny:     true,
	}
	return a
}
func (a *AST) ResetGblParams(params map[string]interface{}) {
	a.globalParams = params
	a.localParams = make(map[string]interface{})
	a.mParams = map[string]interface{}{}
	a.lParams = map[string]interface{}{}
	a.selfParams = map[string]interface{}{}
	a.classParams = map[string]interface{}{}
}

func (a *AST) ParseExpression() types.ExprAST {
	a.depth++ // called depth
	lhs := a.parsePrimary()
	r := a.parseBinOpRHS(0, lhs)
	a.depth--
	if a.depth == 0 && a.currIndex != len(a.Tokens) && a.Err == nil {
		a.Err = errors.New(
			fmt.Sprintf("bad expression, reaching the end or missing the operator\n%s",
				ErrPos("", a.currTok.Offset)))
	}

	return r
}

func (a *AST) getNextToken() Token {
	a.currIndex++
	if a.currIndex < len(a.Tokens) {
		a.currTok = a.Tokens[a.currIndex]
		return a.currTok
	}
	return Token{Type: -1}
}

func (a *AST) getTokPrecedence() int {
	if p, ok := conf.Precedence[a.currTok.Tok]; ok {
		return p
	}
	return -1
}

/*
*
解析float类型
*/
func (a *AST) parseFloat() types.FloatExprAST {
	f64, err := strconv.ParseFloat(a.currTok.Tok, 64)
	if err != nil {
		a.Err = errors.New(
			fmt.Sprintf("%v\nwant '(' or '0-9' but get '%s'\n%s",
				err.Error(),
				a.currTok.Tok,
				ErrPos(a.source, a.currTok.Offset)))
		return types.FloatExprAST{}
	}
	n := types.FloatExprAST{
		Val: f64,
		//Str: a.currTok.Tok,
	}
	a.getNextToken()
	return n
}

/*
*

	int
*/
func (a *AST) parseInt() types.IntExprAST {
	intF, err := strconv.Atoi(a.currTok.Tok)
	if err != nil {
		a.Err = errors.New(
			fmt.Sprintf("%v\nwant '(' or '0-9' but get '%s'\n%s",
				err.Error(),
				a.currTok.Tok,
				ErrPos(a.source, a.currTok.Offset)))
		return types.IntExprAST{}
	}
	n := types.IntExprAST{
		Val: intF,
		//Str: a.currTok.Tok,
	}
	a.getNextToken()
	return n
}

/*
*
*解析函数
 */
func (a *AST) parseFunCallerOrConst() types.ExprAST {
	name := a.currTok.Tok
	a.getNextToken()

	//call func
	if a.currTok.Tok == "(" {
		f := types.FunCallerExprAST{}
		if _, ok := a.defFunc[name]; !ok && a.ignoreFunc[name] == false {
			a.Err = errors.New(
				fmt.Sprintf("function `%s` is undefined\n%s",
					name,
					ErrPos(a.source, a.currTok.Offset)))
			return f
		}

		a.getNextToken()
		exprs := make([]types.ExprAST, 0)
		if a.currTok.Tok == ")" {
			// function call without parameters
			// ignore the process of parameter resolution
		} else {
			exprs = append(exprs, a.ParseExpression())
			//打印错误
			if a.Err != nil {
				return f
			}
			for a.currTok.Tok != ")" && a.getNextToken().Type != -1 {
				if a.currTok.Type == conf.COMMA {
					continue
				}
				if a.Err != nil {
					return f
				}
				exprs = append(exprs, a.ParseExpression())
			}
		}

		def := a.defFunc[name]
		if def.argc >= 0 && len(exprs) != def.argc {
			//fmt.Println(a.source[a.currTok.Offset : a.currTok.Offset+200])
			a.Err = errors.New(
				fmt.Sprintf("wrong way calling function `%s`, parameters want %d but get %d\n%s;token=%d",
					name,
					def.argc,
					len(exprs),
					ErrPos(a.source, a.currTok.Offset),
					exprs,
				),
			)
		}
		a.getNextToken()
		f.Name = name
		f.Arg = exprs
		return f
	}
	// call const
	switch name {
	case conf.SYS_NULL:
		return types.NilExprAST{
			Val: nil,
		}
	case conf.SYS_TRUE:
		return types.BoolExprAST{
			Val: true,
		}
	case conf.SYS_FALSE:
		return types.BoolExprAST{
			Val: false,
		}
	case conf.SYS_ZERO:
		return types.ZeroExprAST{
			Val: "__DSL-ZERO__",
		}
	case conf.SYS_SELF:
		return types.SelfExprAST{
			Val: "__DSL-SELF__",
		}
	default:
		a.Err = errors.New(
			fmt.Sprintf("const `%s` is undefined\n%s",
				name,
				ErrPos(a.source, a.currTok.Offset)))
		return nil
	}

}

/*
*
解析变量
*/
func (a *AST) parseVariable() types.ExprAST {
	if a.currTok.Tok == conf.SYS_ENV {
		expr := types.GlobalExprAST{
			Name: a.currTok.Tok,
			Keys: a.currTok.GblExtra[1:],
		}
		a.getNextToken()
		return expr
	}
	if conf.SysLocalSet[a.currTok.Tok] {
		expr := types.GlobalExprAST{
			Name: a.currTok.Tok,
			Keys: a.currTok.GblExtra,
		}
		a.getNextToken()
		return expr
	}
	if conf.SysSecondSet[a.currTok.Tok] {
		if a.currTok.Tok == conf.SYS_OP || a.currTok.Tok == conf.SYS_OPP {
			expr := types.GlobalExprAST{
				Name: a.currTok.Tok,
				Keys: a.currTok.GblExtra[1:],
			}
			a.getNextToken()
			return expr
		} else {
			expr := types.GlobalExprAST{
				Name: conf.SYS_ENV,
				Keys: a.currTok.GblExtra,
			}
			a.getNextToken()
			return expr
		}

	}

	//业务主数据
	if a.currTok.Tok == conf.SYS_M {
		expr := types.MasterExprAST{
			Name: a.currTok.Tok,
			Keys: a.currTok.GblExtra,
		}
		a.getNextToken()
		return expr
	}

	//业务对照数据
	if a.currTok.Tok == conf.SYS_L {
		expr := types.LookupExprAST{
			Name: a.currTok.Tok,
			Keys: a.currTok.GblExtra,
		}
		a.getNextToken()
		return expr
	}

	//类方法
	if a.currTok.Tok == conf.SYS_C {
		expr := types.ClassExprAST{
			Name: a.currTok.Tok,
			Keys: a.currTok.GblExtra,
		}
		a.getNextToken()
		return expr
	}

	//普通的
	expr := types.VariableExprAST{
		Name: a.currTok.Tok,
		Keys: a.currTok.GblExtra[1:],
	}
	a.getNextToken()
	return expr
}

/*
*
解析全局变量
*/
func (a *AST) parseGlobal() types.ExprAST {
	expr := types.GlobalExprAST{}
	if conf.SysSecondSet[a.currTok.Tok] == true {
		expr.Name = conf.SYS_ENV
		expr.Keys = a.currTok.GblExtra
	} else {
		expr.Name = a.currTok.Tok
		expr.Keys = a.currTok.GblExtra[1:]
	}
	a.getNextToken()
	return expr
}

/*
*
解析变量
*/
func (a *AST) parseField() types.ExprAST {
	expr := types.FieldExprAST{
		Name: a.currTok.Tok,
		Keys: a.currTok.GblExtra,
	}
	a.getNextToken()
	return expr
}

/*
*参数引用
 */
func (a *AST) parseParamsReference() types.ExprAST {
	a.getNextToken()
	return types.NilExprAST{}
	//val, bl := maphelper.JDGet(a.params, strings.Split(a.currTok.Tok, "."))
	//if bl == false {
	//	a.Err = errors.New("{" + a.currTok.Tok + "}" + "解析失败!")
	//	return nil
	//}
	//v := types.AnyExprAST{
	//	Val:  val,
	//	Name: a.currTok.Tok,
	//}
	//
	//a.getNextToken()
	//return v
}

/*
*
解析interface
*/
func (a *AST) parseInterface() types.ExprAST {
	v := types.AnyExprAST{
		Val: a.currTok.Val,
		//Name: a.currTok.Tok,
	}
	a.getNextToken()
	return v
}

/*
*
解析input
*/
func (a *AST) parseInput() types.ExprAST {
	v := types.AnyExprAST{
		Val: a.currTok.Val,
		//Name: a.currTok.Tok,
	}
	a.getNextToken()
	return v
}

/*
*
解析Code
*/
func (a *AST) parseCode() types.ExprAST {
	v := types.AnyExprAST{
		Val: a.currTok.Val,
		//Name: a.currTok.Tok,
	}
	a.getNextToken()
	return v
}

/*
*
解析变量
*/
func (a *AST) parsePrimary() types.ExprAST {
	switch a.currTok.Type {
	case conf.Identifier:
		return a.parseFunCallerOrConst()
	case conf.Literal:
		return a.parseFloat()
	case conf.Operator:
		if a.currTok.Tok == "(" {
			if a.getNextToken().Type == -1 {
				a.Err = errors.New(
					fmt.Sprintf("want '(' or '0-9' but get EOF\n%s",
						ErrPos(a.source, a.currTok.Offset)))
				return nil
			}
			e := a.ParseExpression()
			if e == nil {
				return nil
			}
			if a.currTok.Tok != ")" {
				a.Err = errors.New(
					fmt.Sprintf("want ')' but get %s\n%s",
						a.currTok.Tok,
						ErrPos(a.source, a.currTok.Offset)))
				return nil
			}
			a.getNextToken()
			return e
		} else if a.currTok.Tok == "-" { //右表达式
			if a.getNextToken().Type == -1 {
				a.Err = errors.New(
					fmt.Sprintf("want '0-9' but get '-'\n%s",
						ErrPos(a.source, a.currTok.Offset)))
				return nil
			}
			bin := types.BinaryExprAST{
				Op:  "-",
				Lhs: types.FloatExprAST{},
				Rhs: a.parsePrimary(),
			}
			return bin
		} else if a.currTok.Tok == "!" { //右表达式
			if a.getNextToken().Type == -1 {
				a.Err = errors.New(
					fmt.Sprintf("! next is empty! %s",
						ErrPos(a.source, a.currTok.Offset)))
				return nil
			}
			bin := types.BinaryExprAST{
				Op:  "!",
				Lhs: types.FloatExprAST{},
				Rhs: a.parsePrimary(),
			}
			return bin
		} else {
			return a.parseFloat()
		}
	case conf.COMMA:
		a.Err = errors.New(
			fmt.Sprintf("want '(' or '0-9' but get %s\n%s",
				a.currTok.Tok,
				ErrPos(a.source, a.currTok.Offset)))
		return nil
	case conf.VARIABLE:
		return a.parseVariable()
	case conf.INPUT:
		return a.parseInput()
	case conf.FIELD:
		return a.parseField()
	case conf.GLOBAL:
		return a.parseGlobal()
	default:
		return nil
	}
}

/*
*
解析复合型
*/
func (a *AST) parseBinOpRHS(execPrec int, lhs types.ExprAST) types.ExprAST {
	for {
		tokPrec := a.getTokPrecedence()
		if tokPrec < execPrec {
			return lhs
		}

		//表示为符号
		binOp := a.currTok.Tok
		if a.getNextToken().Type == -1 {
			a.Err = errors.New(
				fmt.Sprintf("want '(' or '0-9' but get EOF\n%s",
					ErrPos(a.source, a.currTok.Offset)))
			return nil
		}

		//解析符号后边的
		rhs := a.parsePrimary()
		if rhs == nil {
			return nil
		}

		nextPrec := a.getTokPrecedence()
		if tokPrec < nextPrec {
			rhs = a.parseBinOpRHS(tokPrec+1, rhs)
			if rhs == nil {
				return nil
			}
		}
		lhs = types.BinaryExprAST{
			Op:  binOp,
			Lhs: lhs,
			Rhs: rhs,
		}
	}
}
