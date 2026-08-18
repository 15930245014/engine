package dsl

import (
	"errors"
	"gitlab.shudieds.com/mec/lib/utils/convert"
	"gitlab.shudieds.com/mec/lib/utils/maphelper"
	"gitlab.shudieds.com/zxh/engine/conf"
	"gitlab.shudieds.com/zxh/engine/types"
	"gitlab.shudieds.com/zxh/engine/utils"
	"strings"
)

func ErrPos(s string, pos int) string {
	r := strings.Repeat("-", len(s)) + "\n"
	s += "\n"
	for i := 0; i < pos; i++ {
		s += " "
	}
	s += "^\n"
	return r + s + r
}

// ExprASTResult is a Top level function
// AST traversal
// if an arithmetic runtime error occurs, a panic exception is thrown
func (a *AST) ExprASTResult(expr types.ExprAST) (interface{}, error) {
	switch expr.(type) {
	case types.BinaryExprAST: //基本运算符
		ast := expr.(types.BinaryExprAST)
		switch ast.Op {
		case "+":
			//l
			l, err := a.ExprASTResult(ast.Lhs)
			if err != nil {
				return nil, err
			}

			//r
			r, err := a.ExprASTResult(ast.Rhs)
			if err != nil {
				return nil, err
			}
			return convert.Add(l, r)
		case "-":
			//l
			l, err := a.ExprASTResult(ast.Lhs)
			if err != nil {
				return nil, err
			}

			//r
			r, err := a.ExprASTResult(ast.Rhs)
			if err != nil {
				return nil, err
			}
			return convert.Sub(l, r)
		case "*":
			//l
			l, err := a.ExprASTResult(ast.Lhs)
			if err != nil {
				return nil, err
			}
			//r
			r, err := a.ExprASTResult(ast.Rhs)
			if err != nil {
				return nil, err
			}
			return convert.Mul(l, r)
		case "/":
			//l
			l, err := a.ExprASTResult(ast.Lhs)
			if err != nil {
				return nil, err
			}
			//r
			r, err := a.ExprASTResult(ast.Rhs)
			if err != nil {
				return nil, err
			}
			return convert.Div(l, r)
		case "%":
			//l
			l, err := a.ExprASTResult(ast.Lhs)
			if err != nil {
				return nil, err
			}

			//r
			r, err := a.ExprASTResult(ast.Rhs)
			if err != nil {
				return nil, err
			}
			return convert.Mod(l, r)
		case "^":
			l, err := a.ExprASTResult(ast.Lhs)
			if err != nil {
				return nil, err
			}

			//r
			r, err := a.ExprASTResult(ast.Rhs)
			if err != nil {
				return nil, err
			}
			return convert.Pow(l, r)
		case ">":
			//l
			l, err := a.ExprASTResult(ast.Lhs)
			if err != nil {
				return nil, err
			}

			//r
			r, err := a.ExprASTResult(ast.Rhs)
			if err != nil {
				return nil, err
			}
			res, ok := convert.Compare(l, r)
			if ok != nil {
				return nil, errors.New("err:> 前后数据类型不一致或不支持的比较类型")
			}
			return res == 1, nil
		case ">=":
			//l
			l, err := a.ExprASTResult(ast.Lhs)
			if err != nil {
				return nil, err
			}

			//r
			r, err := a.ExprASTResult(ast.Rhs)
			if err != nil {
				return nil, err
			}
			res, ok := convert.Compare(l, r)
			if ok != nil {
				return nil, errors.New("err:>= 前后数据类型不一致或不支持的比较类型")
			}
			return res >= 0, nil
		case "<":
			//l
			l, err := a.ExprASTResult(ast.Lhs)
			if err != nil {
				return nil, err
			}

			//r
			r, err := a.ExprASTResult(ast.Rhs)
			if err != nil {
				return nil, err
			}
			res, ok := convert.Compare(l, r)
			if ok != nil {
				return nil, errors.New("err:< 前后数据类型不一致或不支持的比较类型")
			}
			return res == -1, nil
		case "<=":
			//l
			l, err := a.ExprASTResult(ast.Lhs)
			if err != nil {
				return nil, err
			}

			//r
			r, err := a.ExprASTResult(ast.Rhs)
			if err != nil {
				return nil, err
			}
			res, ok := convert.Compare(l, r)
			if ok != nil {
				return nil, errors.New("err:<= 前后数据类型不一致或不支持的比较类型")
			}
			return res <= 0, nil
		case "==":
			//l
			l, err := a.ExprASTResult(ast.Lhs)
			if err != nil {
				return nil, err
			}

			//r
			r, err := a.ExprASTResult(ast.Rhs)
			if err != nil {
				return nil, err
			}
			res, ok := convert.Compare(l, r)
			if ok != nil {
				return nil, errors.New("err:== 前后数据类型不一致或不支持的比较类型")
			}
			return res == 0, nil
		case "!=":
			//l
			l, err := a.ExprASTResult(ast.Lhs)
			if err != nil {
				return nil, err
			}

			//r
			r, err := a.ExprASTResult(ast.Rhs)
			if err != nil {
				return nil, err
			}
			res, ok := convert.Compare(l, r)
			if ok != nil {
				return nil, errors.New("err:!= 前后数据类型不一致或不支持的比较类型")
			}
			return res != 0, nil
		case "&&":
			//l
			l, err := a.ExprASTResult(ast.Lhs)
			if err != nil {
				return nil, err
			}
			bl, ok := convert.Bool(l)
			if ok == false {
				return nil, errors.New("err:&& 数据类型返回应为bool")
			}
			if bl == false {
				return false, nil
			}
			//r
			r, err := a.ExprASTResult(ast.Rhs)
			if err != nil {
				return nil, err
			}

			bl, ok = convert.Bool(r)
			if ok == false {
				return nil, errors.New("err:&& 数据类型返回应为bool")
			}
			return bl, nil
		case "||":
			//l
			l, err := a.ExprASTResult(ast.Lhs)
			if err != nil {
				return nil, err
			}
			bl, ok := convert.Bool(l)
			if ok == false {
				return nil, errors.New("err:|| 数据类型返回应为bool")
			}
			if bl == true {
				return true, nil
			}
			//r
			r, err := a.ExprASTResult(ast.Rhs)
			if err != nil {
				return nil, err
			}

			bl, ok = convert.Bool(r)
			if ok == false {
				return nil, errors.New("err:&& 数据类型返回应为bool")
			}
			return bl, nil
		case "!":
			r, err := a.ExprASTResult(ast.Rhs)
			if err != nil {
				return nil, err
			}
			bl, ok := convert.Bool(r)
			if ok == false {
				return nil, errors.New("err:! 数据类型返回应为bool")
			}
			return !bl, nil
		case "->":
			lv, err := a.ExprASTResult(ast.Lhs)
			if err != nil {
				return nil, err
			}
			a.selfParams["__DSL-SELF__"] = lv
			r, err := a.ExprASTResult(ast.Rhs)
			if err != nil {
				return nil, err
			}
			return r, nil
		default:
			return nil, errors.New("不支持的运算符")
		}
	case types.FunCallerExprAST:
		f := expr.(types.FunCallerExprAST)
		def := a.defFunc[f.Name]
		val, err := def.fun(f.Arg...)
		if err != nil {
			return nil, err
		}
		return val, nil
	case types.VariableExprAST:
		varExpr := expr.(types.VariableExprAST)
		var v types.ExprAST
		//fmt.Println(varExpr)
		err := a.variableCalculate.Set(varExpr.Name, a.globalParams, true, nil)
		if err != nil {
			return nil, err
		}
		v, err = a.variableCalculate.ParseVariable()

		if err != nil {
			return nil, err
		}
		if a.dataType[v.GetType()] == false {
			return nil, errors.New("不合法的数据类型！" + v.GetType())
		}
		if len(varExpr.Keys) >= 0 {
			v2, ok := maphelper.JDGet(v.GetVal(), varExpr.Keys)
			if ok == false {
				return nil, errors.New("路径不合法:$" + varExpr.Name + "." + strings.Join(varExpr.Keys, "."))
			}
			return v2, nil
		} else {
			return v.GetVal(), nil
		}
	case types.FieldExprAST:
		varExpr := expr.(types.FieldExprAST)
		var v types.ExprAST
		err := a.fieldCalculate.Set(varExpr.Name, a.globalParams, true, nil)
		if err != nil {
			return nil, err
		}
		v, err = a.fieldCalculate.ParseField()

		if err != nil {
			return nil, err
		}
		if a.dataType[v.GetType()] == false {
			return nil, errors.New("不合法的数据类型！" + v.GetType())
		}
		if len(varExpr.Keys) >= 0 {
			v2, ok := maphelper.JDGet(v.GetVal(), varExpr.Keys)
			if ok == false {
				return nil, errors.New("路径不合法:@" + varExpr.Name + "." + strings.Join(varExpr.Keys, "."))
			}
			return v2, nil
		} else {
			return v.GetVal(), nil
		}
	case types.GlobalExprAST:
		var v interface{}
		var ok bool
		gblExpr := expr.(types.GlobalExprAST)
		if conf.SysLocalSet[gblExpr.Name] == true {
			v, ok = maphelper.JDGet(a.localParams, gblExpr.Keys)
		} else {
			if gblExpr.Name == conf.SYS_OP {
				op, bl := utils.GetOp(a.globalParams)
				if !bl {
					return nil, errors.New("全局变量OP" + strings.Join(gblExpr.Keys, ".") + "不存在")
				}

				v, ok = maphelper.JDGet(op, gblExpr.Keys)
			} else if gblExpr.Name == conf.SYS_OPP {
				opp, bl := utils.GetOpp(a.globalParams)
				if !bl {
					return nil, errors.New("全局变量OPP" + strings.Join(gblExpr.Keys, ".") + "不存在")
				}
				v, ok = maphelper.JDGet(opp, gblExpr.Keys)

			} else {
				v, ok = maphelper.JDGet(a.globalParams, gblExpr.Keys)
			}
		}
		if ok == false {
			return nil, errors.New("全局变量" + strings.Join(gblExpr.Keys, ".") + "不存在")
		}
		return v, nil
	case types.NilExprAST:
		return expr.GetVal(), nil
	case types.ZeroExprAST:
		return expr.GetVal(), nil
	case types.SelfExprAST:
		return a.selfParams["__DSL-SELF__"], nil
	case types.MasterExprAST:
		master := expr.(types.MasterExprAST)
		if len(master.Keys) < 2 {
			return nil, errors.New("主数据变量M后面必须指定对象类型")
		}
		a.mParams["__M_NAME__"] = master.Keys[1]
		return nil, nil
	case types.LookupExprAST:
		lookup := expr.(types.LookupExprAST)
		if len(lookup.Keys) < 2 {
			return nil, errors.New("对照数据变量L后面必须指定对象类型")
		}
		a.lParams["__L_NAME__"] = lookup.Keys[1]
		return nil, nil
	case types.ClassExprAST:
		class := expr.(types.ClassExprAST)
		if len(class.Keys) < 2 {
			return nil, errors.New("变量C后面必须指定对象类型")
		}
		a.classParams["__C_NAME__"] = class.Keys[1]
		return nil, nil
	default:
		if a.dataType[expr.GetType()] == false {
			return nil, errors.New("不合法的数据类型！" + expr.GetType())
		}
		return expr.GetVal(), nil
	}
}
