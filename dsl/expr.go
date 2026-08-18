package dsl

import (
	"context"
	"errors"
	"gitlab.shudieds.com/mec/lib/utils/convert"
	"gitlab.shudieds.com/mec/lib/utils/uuid"
	_interface "gitlab.shudieds.com/zxh/engine/interface"
	"strings"
)

/**
  公式解析类
*/

type ExprCalculate struct {
	c                 context.Context
	expr              string
	params            map[string]interface{}
	variableCalculate _interface.Variable
	fieldCalculate    _interface.Field
	astCache          map[string]*AST
	uniqueId          string
}

func NewExprCalculate(c context.Context) *ExprCalculate {
	expr := new(ExprCalculate)
	expr.c = c
	expr.astCache = make(map[string]*AST)
	return expr
}
func (e *ExprCalculate) SetC(c context.Context) {
	e.c = c
}

/**
  uniqueId公式唯一标识 避免直接用公式做key造成的大key问题影响性能
*/

func (e *ExprCalculate) Set(s string, uniqueId string, params map[string]interface{}, eName string) {
	e.params = params
	e.expr = strings.TrimSpace(s)

	//e.expr = s
	if len(uniqueId) == 0 {
		e.uniqueId = uuid.GenUuidV3(e.expr)
	} else {
		e.uniqueId = uniqueId
	}
}

func (e *ExprCalculate) SetVariableCalculate(variable _interface.Variable) {
	e.variableCalculate = variable
}
func (e *ExprCalculate) SetFieldCalculate(field _interface.Field) {
	e.fieldCalculate = field
}

/**
  检查输入

*/

func (e *ExprCalculate) CheckFormula() (bool, error) {
	expr := e.expr
	if len(expr) == 0 {
		return false, errors.New("formal is empty!")
	}
	tokens, err := GetTokens(expr)

	if err != nil {
		return false, err
	}

	ast := NewAST(e.c, tokens, expr, nil, e.variableCalculate, e.fieldCalculate)

	//ast.SetIsCheck(true)
	if ast.Err != nil {
		return false, ast.Err
	}

	//builder
	_ = ast.ParseExpression()
	if ast.Err != nil {
		return false, ast.Err
	}

	return true, nil
}

func (e *ExprCalculate) SetAstCache(cache interface{}) {
	astC, ok := cache.(map[string]*AST)
	if ok {
		e.astCache = astC
	}
}
func (e *ExprCalculate) GetAstCache() interface{} {
	return e.astCache
}

//(1+2 - 5)
/**
  计算
*/
func (e *ExprCalculate) Calculate() (interface{}, error) {
	//判断是否为空
	if len(e.expr) == 0 {
		return nil, nil
	}
	//expr, err := e.ConvertExpr()
	//if err != nil {
	//	return nil, err
	//}
	//fmt.Println("expr", expr)
	if e.astCache[e.uniqueId] == nil {
		tokens, err := GetTokens(e.expr)
		//for i := 0; i < len(tokens); i++ {
		//	fmt.Println(tokens[i].ToString())
		//}
		if err != nil {
			return nil, err
		}

		ast := NewAST(e.c, tokens, e.expr, e.params, e.variableCalculate, e.fieldCalculate)
		if ast.Err != nil {
			return nil, ast.Err
		}

		//builder
		ast.root = ast.ParseExpression()
		if ast.Err != nil {
			return nil, ast.Err
		}
		e.astCache[e.uniqueId] = ast
	} else {
		e.astCache[e.uniqueId].ResetGblParams(e.params)
	}

	//ast expr
	r, err := e.astCache[e.uniqueId].ExprASTResult(e.astCache[e.uniqueId].root)

	return r, err
}

/**
清除缓存
*/

func (e *ExprCalculate) ClearCache() {
	//e.astCache = make(map[string]*AST)
	e.params = make(map[string]interface{})

}

/*
类型转化
*/

func (e *ExprCalculate) ConvertExpr() (string, error) {
	//todo: or(1,$b->and())->or()

	res, err := e.dfs(e.expr, 0, len(e.expr)-1)
	if err != nil {
		return "", err
	}
	return convert.BytesToString(res), nil
}

func (f *ExprCalculate) dfs(expr string, s, e int) ([]byte, error) {
	//操作符
	operator := make(map[byte]bool)
	operator['+'] = true
	operator['-'] = true
	operator['*'] = true
	operator['/'] = true
	operator['^'] = true
	operator['%'] = true
	operator[','] = true

	operator['>'] = true
	operator['<'] = true
	operator['='] = true

	operator['&'] = true
	operator['|'] = true
	operator['!'] = true

	//返回值
	var res []byte
	var pre []byte //前驱
	i := s
	for i <= e {
		//判断是否为运算符
		if expr[i] == '-' && i+1 < len(expr) && expr[i+1] == '>' {
			//i -> '>'
			i++

			//获取函数名
			j := i + 1 //>后边的第一位
			for j < len(expr) && expr[j] != '(' {
				j++
			}

			//此时j为 '('
			if j >= len(expr) || j-i == 1 {
				return nil, errors.New("函数表达式错误")
			}

			//获取函数名称
			defName := expr[i+1 : j]

			//判断
			if tokenCheckDef[defName] == false {
				return nil, errors.New(defName + "函数名称不合法!")
			}

			//获取函数参数
			start := j + 1 //括号内第一个元素
			tip := -1
			j = start
			for j < len(expr) {
				if expr[j] == '(' {
					tip--
				}
				if expr[j] == ')' {
					tip++
				}
				//判断是否为最右边的括号
				if tip == 0 {
					break
				}
				j++
			}
			if j >= len(expr) {
				return nil, errors.New(defName + "无效的()")
			}

			//判断前缀参数
			if len(pre) == 0 {
				return nil, errors.New("->前边参数error！")
			}

			//
			next := []byte{}

			//j此时为最右边的括号
			if j == start {
				next = append(next, []byte(defName)...)
				next = append(next, '(')
				next = append(next, pre...)
				next = append(next, ')')
			} else {
				rt, err := f.dfs(expr, start, j-1)
				if err != nil {
					return nil, err
				}
				next = append(next, []byte(defName)...)
				next = append(next, '(')
				next = append(next, pre...)
				next = append(next, ',')
				next = append(next, rt...)
				next = append(next, ')')
			}
			pre = next

			//从 j + 1开始
			i = j + 1
			continue
		}

		//区分参数中的 +-*/等
		if expr[i] == '"' {
			pre = append(pre, '"')
			j := i + 1
			for j < len(expr) {
				pre = append(pre, expr[j])
				if expr[j] == '"' && expr[j-1] != '\\' {
					break
				}
				j++
			}

			if j >= len(expr) {
				return nil, errors.New("无效的字符文本")
			}
			i = j + 1
			continue
		}

		//否则
		if operator[expr[i]] == true {
			res = append(res, pre...)
			res = append(res, expr[i])
			pre = []byte{}
		} else {
			pre = append(pre, expr[i])
		}
		//继续遍历
		i++
	}
	//判断最后
	if len(pre) > 0 {
		res = append(res, pre...)
	}
	return res, nil
}

/*
*
is_char 是否为字母
*/
func (e *ExprCalculate) isChar(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}
