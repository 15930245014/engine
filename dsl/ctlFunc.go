package dsl

import (
	"errors"
	"github.com/shopspring/decimal"
	"gitlab.shudieds.com/mec/lib/consts"
	"gitlab.shudieds.com/mec/lib/utils/convert"
	"gitlab.shudieds.com/zxh/engine/conf"
	"gitlab.shudieds.com/zxh/engine/types"
	eUtils "gitlab.shudieds.com/zxh/engine/utils"
	"strconv"
	"time"
)

func (a *AST) defIf(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	bl, _ := convert.Bool(v)
	if bl {
		return a.ExprASTResult(expr[1])
	}
	return a.ExprASTResult(expr[2])
}

/**
case函数
*/

func (a *AST) defCase(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"case": v,
		"do":   expr[1],
	}, nil
}

/*
* 默认值查询
* default函数
 */
func (a *AST) defDefault(expr ...types.ExprAST) (interface{}, error) {
	return map[string]interface{}{
		"def":  true,
		"case": true,
		"do":   expr[0],
	}, nil
}

/*
*

	switch
*/
func (a *AST) defSwitch(expr ...types.ExprAST) (interface{}, error) {
	if len(expr) == 0 {
		return nil, errors.New("defSwitch 参数不能为空!")
	}
	//判断长度
	if len(expr) <= 1 {
		return nil, errors.New("defSwitch 参数应该大于1个！")
	}

	where, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	for i := 1; i < len(expr); i++ {
		v, err := a.ExprASTResult(expr[i])
		if err != nil {
			return nil, err
		}

		res := v.(map[string]interface{})
		//default
		if _, ok := res["def"]; ok {
			return a.ExprASTResult(res["do"].(types.ExprAST))
		}
		com, _ := convert.Compare(res["case"], where)
		if com == 0 {
			return a.ExprASTResult(res["do"].(types.ExprAST))
		}
	}

	return nil, nil
}

/*
  - item
    index 会引用
    rt
    do("continue/break/返回值","sun(RT,item.s4)"),do2

    Choice(case(),case())
*/
func (a *AST) defChoice(expr ...types.ExprAST) (interface{}, error) {
	if len(expr) == 0 {
		return nil, errors.New("defChoice 参数不能为空!")
	}

	for i := 0; i < len(expr); i++ {
		v, err := a.ExprASTResult(expr[i])
		if err != nil {
			return nil, err
		}
		res := v.(map[string]interface{})
		bl, ok := convert.Bool(res["case"])
		if ok == true && bl == true {
			return a.ExprASTResult(res["do"].(types.ExprAST))
		}
	}
	//返回值均为nil
	return nil, nil
}

/*
*
defForeach <script>
var sum = 0;
var numbers = [65, 44, 12, 4];

	function myFunction(item) {
	    sum += item;
	    demo.innerHTML = sum;
	}
*/
//func (a *AST) defForeach(exprs ...types.ExprAST) (interface{}, error) {
//	//判断参数
//	if len(exprs) < 2 {
//		return nil, errors.New("defForeach 参数数目错误")
//	}
//	arr, err := a.ExprASTResult(exprs[0])
//	if err != nil {
//		return nil, err
//	}
//	//判断数组类型
//	t, bl := eUtils.ArrType(arr, 1)
//	if bl == false {
//		return nil, errors.New("defForeach 错误的参数类型")
//	}
//	//声明返回值和执行表达式
//	rt := make([]interface{}, len(exprs))
//	keys := make(map[int]string)
//	doExpr := make([]map[string]interface{}, len(exprs))
//	result := make(map[string]interface{})
//
//	//哨兵机制 避免每次循环都去判断
//	for i := 1; i < len(exprs); i++ {
//		v, err := a.ExprASTResult(exprs[i])
//		if err != nil {
//			return result, nil
//		}
//
//		res := v.(map[string]interface{})
//		keys[i] = res["key"].(string)
//		rType := res["rType"].(string)
//		switch rType {
//		case conf.ReturnInt:
//			doV, err := a.ExprASTResult(res["init"].(types.ExprAST))
//			if err != nil {
//				return rt, err
//			}
//			if doV != nil {
//				doRt, ok := convert.Int(doV)
//				if ok == false {
//					return rt, errors.New("each第" + strconv.Itoa(i) + "参数 计算结果无法转变成int!")
//				}
//				rt[i] = doRt
//			} else {
//				rt[i] = 0
//			}
//		case conf.ReturnFloat:
//			doV, err := a.ExprASTResult(res["init"].(types.ExprAST))
//			if err != nil {
//				return rt, err
//			}
//			if doV != nil {
//				doRt, ok := convert.Float64(doV)
//				if ok == false {
//					return rt, errors.New("each第" + strconv.Itoa(i) + "参数 计算结果无法转变成float")
//				}
//				rt[i] = doRt
//			} else {
//				rt[i] = float64(0)
//			}
//		case conf.ReturnStr:
//			doV, err := a.ExprASTResult(res["init"].(types.ExprAST))
//			if err != nil {
//				return rt, err
//			}
//			if rt != nil {
//				doRt, ok := convert.String(doV)
//				if ok == false {
//					return rt, errors.New("each第" + strconv.Itoa(i) + "参数 计算结果无法转变成str")
//				}
//				rt[i] = doRt
//			} else {
//				rt[i] = ""
//			}
//		case conf.ReturnBool:
//			doV, err := a.ExprASTResult(res["init"].(types.ExprAST))
//			if err != nil {
//				return rt, err
//			}
//			if doV != nil {
//				doRt, ok := convert.Bool(doV)
//				if ok == false {
//					return rt, errors.New("each第" + strconv.Itoa(i) + "参数 计算结果无法转变成bool")
//				}
//				rt[i] = doRt
//			} else {
//				rt[i] = false
//			}
//		case conf.ReturnDate:
//			doV, err := a.ExprASTResult(res["init"].(types.ExprAST))
//			if err != nil {
//				return rt, err
//			}
//			if doV != nil {
//				doRt, ok := convert.Date(doV)
//				if ok == false {
//					return rt, errors.New("each第" + strconv.Itoa(i) + "参数 计算结果无法转变成date类型!")
//				}
//				rt[i] = doRt
//			} else {
//				//初始化为空的时间
//				rt[i] = time.Time{}
//			}
//		case conf.ReturnDecimal:
//			doV, err := a.ExprASTResult(res["init"].(types.ExprAST))
//			if err != nil {
//				return rt, err
//			}
//			if doV != nil {
//				doRt, ok := convert.Decimal(doV)
//				if ok == false {
//					return rt, errors.New("each第" + strconv.Itoa(i) + "参数 计算结果无法转变成返回decimal!")
//				}
//				rt[i] = doRt
//			} else {
//				rt[i] = decimal.NewFromInt(int64(0))
//			}
//		case conf.ReturnAny:
//			doV, err := a.ExprASTResult(res["init"].(types.ExprAST))
//			if err != nil {
//				return rt, err
//			}
//			if doV != nil {
//				rt[i] = doV
//			} else {
//				rt[i] = nil
//			}
//		case conf.ReturnMap:
//			doV, err := a.ExprASTResult(res["init"].(types.ExprAST))
//			if err != nil {
//				return rt, err
//			}
//			if doV != nil {
//				vv, ok := convert.ToMap(doV)
//				if ok == false {
//					return rt, errors.New("each第" + strconv.Itoa(i) + "参数 计算结果无法转变成map!")
//				}
//				rt[i] = vv
//			} else {
//				rt[i] = make(map[string]interface{})
//			}
//		case conf.ReturnArray:
//			doV, err := a.ExprASTResult(res["init"].(types.ExprAST))
//			if err != nil {
//				return rt, err
//			}
//			if doV != nil {
//				_, ok := eUtils.ArrType(doV, 1)
//				if ok == false {
//					return rt, errors.New("each第" + strconv.Itoa(i) + "参数 计算结果无法转变成arr!")
//				}
//				rt[i] = doV
//
//			} else {
//				rt[i] = []interface{}{}
//			}
//		case conf.ReturnContinue:
//			doV, err := a.ExprASTResult(res["init"].(types.ExprAST))
//			if err != nil {
//				return rt, err
//			}
//			if doV != nil {
//				rt[i] = doV
//			}
//		case conf.ReturnBreak:
//			doV, err := a.ExprASTResult(res["init"].(types.ExprAST))
//			if err != nil {
//				return rt, err
//			}
//			if doV != nil {
//				rt[i] = doV
//			}
//		case conf.ReturnArrStr:
//			doV, err := a.ExprASTResult(res["init"].(types.ExprAST))
//			if err != nil {
//				return rt, err
//			}
//			if doV != nil {
//				_, ok := doV.([]string)
//				if ok == false {
//					return rt, errors.New("each第" + strconv.Itoa(i) + "参数 计算结果无法转arr.str")
//				}
//				rt[i] = doV
//			} else {
//				rt[i] = []string{}
//			}
//		default:
//			return rt, errors.New("each第" + strconv.Itoa(i) + "参数:不支持的数据返回类型")
//		}
//
//		//写入到任务
//		doExpr[i] = res
//	}
//	//判断数组类别
//	switch t {
//	case conf.ReturnMap:
//		//生成对象数组
//		arrObj, _ := convert.ToArrMap(arr)
//		if len(arrObj) > 0 {
//			//执行结果
//			isBreak := false
//			for i := 0; i < len(arrObj) && false == isBreak; i++ {
//				a.localParams[conf.SYS_E] = arrObj[i]
//				a.localParams[conf.SYS_I] = i
//				isContinue := false
//				for j := 1; j < len(doExpr) && isContinue == false; j++ {
//					a.localParams[conf.SYS_V] = rt[j]
//					doV, err := a.ExprASTResult(doExpr[j]["do"].(types.ExprAST))
//					if err != nil {
//						return nil, err
//					}
//					rType := doExpr[j]["rType"].(string)
//					switch rType {
//					case conf.ReturnInt:
//						doRt, ok := convert.Int(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变int")
//						}
//						rt[j] = doRt
//					case conf.ReturnFloat:
//						doRt, ok := convert.Float64(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变float")
//						}
//						rt[j] = doRt
//					case conf.ReturnStr:
//						doRt, ok := convert.String(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变str")
//						}
//						rt[j] = doRt
//					case conf.ReturnBool:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成bool")
//						}
//						rt[j] = doRt
//					case conf.ReturnDate:
//						doRt, ok := convert.Date(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成date")
//						}
//						rt[j] = doRt
//					case conf.ReturnDecimal:
//						doRt, ok := convert.Decimal(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成decimal")
//						}
//						rt[j] = doRt
//					case conf.ReturnAny:
//						rt[j] = doV
//					case conf.ReturnMap:
//						doRt, ok := doV.(map[string]interface{})
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成map")
//						}
//						rt[j] = doRt
//
//					case conf.ReturnArray:
//						doRt, ok := eUtils.ArrType(doV, 1)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成arr!")
//						}
//						rt[j] = doRt
//					case conf.ReturnContinue:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成continue!")
//						}
//						if doRt == true {
//							isContinue = true
//						}
//
//					case conf.ReturnBreak:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成break!")
//						}
//						if doRt == true {
//							isBreak = true
//							isContinue = true
//						}
//					case conf.ReturnArrStr:
//						doRt, ok := doV.([]string)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成arr.str!")
//						}
//						rt[j] = doRt
//					default:
//						return rt, errors.New("each第" + strconv.Itoa(i) + "参数:不支持的数据返回类型")
//					}
//				}
//			}
//
//		}
//	case conf.ReturnInt:
//		arrInt, _ := arr.([]int)
//		if len(arrInt) > 0 {
//			//执行结果
//			isBreak := false
//			for i := 0; i < len(arrInt) && isBreak == false; i++ {
//				a.localParams[conf.SYS_E] = arrInt[i]
//				a.localParams[conf.SYS_I] = i
//				isContinue := false
//				for j := 1; j < len(doExpr) && isContinue == false; j++ {
//					a.localParams[conf.SYS_V] = rt[j]
//					doV, err := a.ExprASTResult(doExpr[j]["do"].(types.ExprAST))
//					if err != nil {
//						return nil, err
//					}
//					rType := doExpr[j]["rType"].(string)
//					switch rType {
//					case conf.ReturnInt:
//						doRt, ok := convert.Int(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成int!")
//						}
//						rt[j] = doRt
//					case conf.ReturnFloat:
//						doRt, ok := convert.Float64(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成float!")
//						}
//						rt[i] = doRt
//					case conf.ReturnStr:
//						doRt, ok := convert.String(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返str!")
//						}
//						rt[i] = doRt
//					case conf.ReturnBool:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返bool")
//						}
//						rt[i] = doRt
//					case conf.ReturnDate:
//						doRt, ok := convert.Date(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成date!")
//						}
//						rt[i] = doRt
//					case conf.ReturnDecimal:
//						doRt, ok := convert.Decimal(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成decimal")
//						}
//						rt[i] = doRt
//					case conf.ReturnAny:
//						rt[i] = doV
//					case conf.ReturnMap:
//						doRt, ok := doV.(map[string]interface{})
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成map!")
//						}
//						rt[i] = doRt
//
//					case conf.ReturnArray:
//						doRt, ok := eUtils.ArrType(doV, 1)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成arr!")
//						}
//						rt[i] = doRt
//					case conf.ReturnContinue:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成continue!")
//						}
//						if doRt == true {
//							isContinue = true
//						}
//
//					case conf.ReturnBreak:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返break!")
//						}
//						if doRt == true {
//							isBreak = true
//							//哨兵减少本轮循环的条件增加isBreak ==false
//							isContinue = true
//						}
//					case conf.ReturnArrStr:
//						doRt, ok := doV.([]string)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返arr.str!")
//						}
//						rt[i] = doRt
//					default:
//						return rt, errors.New("each第" + strconv.Itoa(i) + "参数:不支持的数据返回类型")
//					}
//				}
//			}
//		}
//	case conf.ReturnFloat:
//		arrFloat, _ := arr.([]float64)
//		if len(arrFloat) > 0 {
//			//执行结果
//			isBreak := false
//			for i := 0; i < len(arrFloat) && isBreak == false; i++ {
//				a.localParams[conf.SYS_E] = arrFloat[i]
//				a.localParams[conf.SYS_I] = i
//				isContinue := false
//				for j := 1; j < len(doExpr) && isContinue == false; j++ {
//					a.localParams[conf.SYS_V] = rt[j]
//					doV, err := a.ExprASTResult(doExpr[j]["do"].(types.ExprAST))
//					if err != nil {
//						return nil, err
//					}
//					rType := doExpr[j]["rType"].(string)
//					switch rType {
//					case conf.ReturnInt:
//						doRt, ok := convert.Int(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnFloat:
//						doRt, ok := convert.Float64(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnStr:
//						doRt, ok := convert.String(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnBool:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnDate:
//						doRt, ok := convert.Date(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnDecimal:
//						doRt, ok := convert.Decimal(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnAny:
//						rt[i] = doV
//					case conf.ReturnMap:
//						doRt, ok := doV.(map[string]interface{})
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//
//					case conf.ReturnArray:
//						doRt, ok := eUtils.ArrType(doV, 1)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnContinue:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						if doRt == true {
//							isContinue = true
//						}
//
//					case conf.ReturnBreak:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						if doRt == true {
//							isBreak = true
//							isContinue = true
//						}
//					case conf.ReturnArrStr:
//						doRt, ok := doV.([]string)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					default:
//						return rt, errors.New("each第" + strconv.Itoa(i) + "参数:不支持的数据返回类型")
//					}
//				}
//			}
//		}
//	case conf.ReturnStr:
//		arrStr, _ := arr.([]string)
//		if len(arrStr) > 0 {
//			//执行结果
//			isBreak := false
//			for i := 0; i < len(arrStr) && isBreak == false; i++ {
//				a.localParams[conf.SYS_E] = arrStr[i]
//				a.localParams[conf.SYS_I] = i
//				isContinue := false
//				for j := 1; j < len(doExpr) && isContinue == false; j++ {
//					a.localParams[conf.SYS_V] = rt[j]
//					doV, err := a.ExprASTResult(doExpr[j]["do"].(types.ExprAST))
//					if err != nil {
//						return nil, err
//					}
//					rType := doExpr[j]["rType"].(string)
//					switch rType {
//					case conf.ReturnInt:
//						doRt, ok := convert.Int(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnFloat:
//						doRt, ok := convert.Float64(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnStr:
//						doRt, ok := convert.String(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnBool:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnDate:
//						doRt, ok := convert.Date(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnDecimal:
//						doRt, ok := convert.Decimal(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnAny:
//						rt[j] = doV
//					case conf.ReturnMap:
//						doRt, ok := doV.(map[string]interface{})
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//
//					case conf.ReturnArray:
//						doRt, ok := eUtils.ArrType(doV, 1)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnContinue:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						if doRt == true {
//							isContinue = true
//						}
//
//					case conf.ReturnBreak:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						if doRt == true {
//							isBreak = true
//							isContinue = true
//						}
//					case conf.ReturnArrStr:
//						doRt, ok := doV.([]string)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					default:
//						return rt, errors.New("each第" + strconv.Itoa(i) + "参数:不支持的数据返回类型")
//					}
//				}
//
//			}
//		}
//	case conf.ReturnBool:
//		arrBool, _ := arr.([]bool)
//		if len(arrBool) > 0 {
//			//执行结果
//			isBreak := false
//			for i := 0; i < len(arrBool) && isBreak == false; i++ {
//				a.localParams[conf.SYS_E] = arrBool[i]
//				a.localParams[conf.SYS_I] = i
//				isContinue := false
//				for j := 1; j < len(doExpr) && isContinue == false; j++ {
//					a.localParams[conf.SYS_V] = rt[j]
//					doV, err := a.ExprASTResult(doExpr[j]["do"].(types.ExprAST))
//					if err != nil {
//						return nil, err
//					}
//					rType := doExpr[j]["rType"].(string)
//					switch rType {
//					case conf.ReturnInt:
//						doRt, ok := convert.Int(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnFloat:
//						doRt, ok := convert.Float64(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnStr:
//						doRt, ok := convert.String(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnBool:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnDate:
//						doRt, ok := convert.Date(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnDecimal:
//						doRt, ok := convert.Decimal(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnAny:
//						rt[i] = doV
//					case conf.ReturnMap:
//						doRt, ok := doV.(map[string]interface{})
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//
//					case conf.ReturnArray:
//						doRt, ok := eUtils.ArrType(doV, 1)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnContinue:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						if doRt == true {
//							isContinue = true
//						}
//
//					case conf.ReturnBreak:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						if doRt == true {
//							isBreak = true
//							isContinue = true
//						}
//					case conf.ReturnArrStr:
//						doRt, ok := doV.([]string)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					default:
//						return rt, errors.New("each第" + strconv.Itoa(i) + "参数:不支持的数据返回类型")
//					}
//				}
//
//			}
//		}
//	case conf.ReturnDate:
//		arrDate, _ := arr.([]time.Time)
//		if len(arrDate) > 0 {
//			//执行结果
//			isBreak := false
//			for i := 0; i < len(arrDate) && isBreak == false; i++ {
//				a.localParams[conf.SYS_E] = arrDate[i]
//				a.localParams[conf.SYS_I] = i
//				isContinue := false
//				for j := 1; j < len(doExpr) && isContinue == false; j++ {
//					a.localParams[conf.SYS_V] = rt[j]
//					doV, err := a.ExprASTResult(doExpr[j]["do"].(types.ExprAST))
//					if err != nil {
//						return nil, err
//					}
//					rType := doExpr[j]["rType"].(string)
//					switch rType {
//					case conf.ReturnInt:
//						doRt, ok := convert.Int(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnFloat:
//						doRt, ok := convert.Float64(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnStr:
//						doRt, ok := convert.String(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnBool:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnDate:
//						doRt, ok := convert.Date(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnDecimal:
//						doRt, ok := convert.Decimal(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnAny:
//						rt[i] = doV
//					case conf.ReturnMap:
//						doRt, ok := doV.(map[string]interface{})
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//
//					case conf.ReturnArray:
//						doRt, ok := eUtils.ArrType(doV, 1)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnContinue:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						if doRt == true {
//							isContinue = true
//						}
//
//					case conf.ReturnBreak:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						if doRt == true {
//							isBreak = true
//							isContinue = true
//						}
//					case conf.ReturnArrStr:
//						doRt, ok := doV.([]string)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					default:
//						return rt, errors.New("each第" + strconv.Itoa(i) + "参数:不支持的数据返回类型")
//					}
//				}
//
//			}
//		}
//	case conf.ReturnAny:
//		arrAny, _ := arr.([]interface{})
//		if len(arrAny) > 0 {
//			//执行结果
//			isBreak := false
//			for i := 0; i < len(arrAny) && isBreak == false; i++ {
//				a.localParams[conf.SYS_E] = arrAny[i]
//				a.localParams[conf.SYS_I] = i
//				isContinue := false
//				for j := 1; j < len(doExpr) && isContinue == false; j++ {
//					a.localParams[conf.SYS_V] = rt[j]
//					doV, err := a.ExprASTResult(doExpr[j]["do"].(types.ExprAST))
//					if err != nil {
//						return nil, err
//					}
//					rType := doExpr[j]["rType"].(string)
//					switch rType {
//					case conf.ReturnInt:
//						doRt, ok := convert.Int(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnFloat:
//						doRt, ok := convert.Float64(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnStr:
//						doRt, ok := convert.String(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnBool:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnDate:
//						doRt, ok := convert.Date(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnDecimal:
//						doRt, ok := convert.Decimal(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnAny:
//						rt[i] = doV
//					case conf.ReturnMap:
//						doRt, ok := doV.(map[string]interface{})
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//
//					case conf.ReturnArray:
//						doRt, ok := eUtils.ArrType(doV, 1)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnContinue:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						if doRt == true {
//							isContinue = true
//						}
//
//					case conf.ReturnBreak:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						if doRt == true {
//							isBreak = true
//							isContinue = true
//						}
//					case conf.ReturnArrStr:
//						doRt, ok := doV.([]string)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					default:
//						return rt, errors.New("each第" + strconv.Itoa(i) + "参数:不支持的数据返回类型")
//					}
//				}
//
//			}
//		}
//	case conf.ReturnArray:
//		arrArr, _ := arr.([]interface{})
//		if len(arrArr) > 0 {
//			//执行结果
//			isBreak := false
//			for i := 0; i < len(arrArr) && isBreak == false; i++ {
//				a.localParams[conf.SYS_E] = arrArr[i]
//				a.localParams[conf.SYS_I] = i
//				isContinue := false
//				for j := 1; j < len(doExpr) && isContinue == false; j++ {
//					a.localParams[conf.SYS_V] = rt[j]
//					doV, err := a.ExprASTResult(doExpr[j]["do"].(types.ExprAST))
//					if err != nil {
//						return nil, err
//					}
//					rType := doExpr[j]["rType"].(string)
//					switch rType {
//					case conf.ReturnInt:
//						doRt, ok := convert.Int(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnFloat:
//						doRt, ok := convert.Float64(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnStr:
//						doRt, ok := convert.String(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnBool:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnDate:
//						doRt, ok := convert.Date(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnDecimal:
//						doRt, ok := convert.Decimal(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnAny:
//						rt[i] = doV
//					case conf.ReturnMap:
//						doRt, ok := doV.(map[string]interface{})
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//
//					case conf.ReturnArray:
//						doRt, ok := eUtils.ArrType(doV, 1)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnContinue:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						if doRt == true {
//							isContinue = true
//						}
//
//					case conf.ReturnBreak:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						if doRt == true {
//							isBreak = true
//							isContinue = true
//						}
//					case conf.ReturnArrStr:
//						doRt, ok := doV.([]string)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					default:
//						return rt, errors.New("each第" + strconv.Itoa(i) + "参数:不支持的数据返回类型")
//					}
//				}
//
//			}
//		}
//	case conf.ReturnFloat32:
//		arrArr, _ := arr.([]float32)
//		if len(arrArr) > 0 {
//			//执行结果
//			isBreak := false
//			for i := 0; i < len(arrArr) && isBreak == false; i++ {
//				a.localParams[conf.SYS_E] = float64(arrArr[i])
//				a.localParams[conf.SYS_I] = i
//				isContinue := false
//				for j := 1; j < len(doExpr) && isContinue == false; j++ {
//					a.localParams[conf.SYS_V] = rt[j]
//					doV, err := a.ExprASTResult(doExpr[j]["do"].(types.ExprAST))
//					if err != nil {
//						return nil, err
//					}
//					rType := doExpr[j]["rType"].(string)
//					switch rType {
//					case conf.ReturnInt:
//						doRt, ok := convert.Int(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnFloat:
//						doRt, ok := convert.Float64(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnStr:
//						doRt, ok := convert.String(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnBool:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnDate:
//						doRt, ok := convert.Date(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnDecimal:
//						doRt, ok := convert.Decimal(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnAny:
//						rt[i] = doV
//					case conf.ReturnMap:
//						doRt, ok := doV.(map[string]interface{})
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//
//					case conf.ReturnArray:
//						doRt, ok := eUtils.ArrType(doV, 1)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnContinue:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						if doRt == true {
//							isContinue = true
//						}
//
//					case conf.ReturnBreak:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						if doRt == true {
//							isBreak = true
//							isContinue = true
//						}
//					case conf.ReturnArrStr:
//						doRt, ok := doV.([]string)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					default:
//						return rt, errors.New("each第" + strconv.Itoa(i) + "参数:不支持的数据返回类型")
//					}
//				}
//
//			}
//		}
//	case conf.ReturnInt8:
//		arrArr, _ := arr.([]int8)
//		if len(arrArr) > 0 {
//			//执行结果
//			isBreak := false
//			for i := 0; i < len(arrArr) && isBreak == false; i++ {
//				a.localParams[conf.SYS_E] = int(arrArr[i])
//				a.localParams[conf.SYS_I] = i
//				isContinue := false
//				for j := 1; j < len(doExpr) && isContinue == false; j++ {
//					a.localParams[conf.SYS_V] = rt[j]
//					doV, err := a.ExprASTResult(doExpr[j]["do"].(types.ExprAST))
//					if err != nil {
//						return nil, err
//					}
//					rType := doExpr[j]["rType"].(string)
//					switch rType {
//					case conf.ReturnInt:
//						doRt, ok := convert.Int(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnFloat:
//						doRt, ok := convert.Float64(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnStr:
//						doRt, ok := convert.String(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnBool:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnDate:
//						doRt, ok := convert.Date(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnDecimal:
//						doRt, ok := convert.Decimal(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnAny:
//						rt[i] = doV
//					case conf.ReturnMap:
//						doRt, ok := doV.(map[string]interface{})
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//
//					case conf.ReturnArray:
//						doRt, ok := eUtils.ArrType(doV, 1)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnContinue:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						if doRt == true {
//							isContinue = true
//						}
//
//					case conf.ReturnBreak:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						if doRt == true {
//							isBreak = true
//							isContinue = true
//						}
//					case conf.ReturnArrStr:
//						doRt, ok := doV.([]string)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					default:
//						return rt, errors.New("each第" + strconv.Itoa(i) + "参数:不支持的数据返回类型")
//					}
//				}
//
//			}
//		}
//	case conf.ReturnInt16:
//		arrArr, _ := arr.([]int16)
//		if len(arrArr) > 0 {
//			//执行结果
//			isBreak := false
//			for i := 0; i < len(arrArr) && isBreak == false; i++ {
//				a.localParams[conf.SYS_E] = int(arrArr[i])
//				a.localParams[conf.SYS_I] = i
//				isContinue := false
//				for j := 1; j < len(doExpr) && isContinue == false; j++ {
//					a.localParams[conf.SYS_V] = rt[j]
//					doV, err := a.ExprASTResult(doExpr[j]["do"].(types.ExprAST))
//					if err != nil {
//						return nil, err
//					}
//					rType := doExpr[j]["rType"].(string)
//					switch rType {
//					case conf.ReturnInt:
//						doRt, ok := convert.Int(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnFloat:
//						doRt, ok := convert.Float64(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnStr:
//						doRt, ok := convert.String(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnBool:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnDate:
//						doRt, ok := convert.Date(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnDecimal:
//						doRt, ok := convert.Decimal(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnAny:
//						rt[i] = doV
//					case conf.ReturnMap:
//						doRt, ok := doV.(map[string]interface{})
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//
//					case conf.ReturnArray:
//						doRt, ok := eUtils.ArrType(doV, 1)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnContinue:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						if doRt == true {
//							isContinue = true
//						}
//
//					case conf.ReturnBreak:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						if doRt == true {
//							isBreak = true
//							isContinue = true
//						}
//					case conf.ReturnArrStr:
//						doRt, ok := doV.([]string)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					default:
//						return rt, errors.New("each第" + strconv.Itoa(i) + "参数:不支持的数据返回类型")
//					}
//				}
//
//			}
//		}
//	case conf.ReturnInt32:
//		arrArr, _ := arr.([]int32)
//		if len(arrArr) > 0 {
//			//执行结果
//			isBreak := false
//			for i := 0; i < len(arrArr) && isBreak == false; i++ {
//				a.localParams[conf.SYS_E] = int(arrArr[i])
//				a.localParams[conf.SYS_I] = i
//				isContinue := false
//				for j := 1; j < len(doExpr) && isContinue == false; j++ {
//					a.localParams[conf.SYS_V] = rt[j]
//					doV, err := a.ExprASTResult(doExpr[j]["do"].(types.ExprAST))
//					if err != nil {
//						return nil, err
//					}
//					rType := doExpr[j]["rType"].(string)
//					switch rType {
//					case conf.ReturnInt:
//						doRt, ok := convert.Int(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnFloat:
//						doRt, ok := convert.Float64(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnStr:
//						doRt, ok := convert.String(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnBool:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnDate:
//						doRt, ok := convert.Date(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnDecimal:
//						doRt, ok := convert.Decimal(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnAny:
//						rt[i] = doV
//					case conf.ReturnMap:
//						doRt, ok := doV.(map[string]interface{})
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//
//					case conf.ReturnArray:
//						doRt, ok := eUtils.ArrType(doV, 1)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnContinue:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						if doRt == true {
//							isContinue = true
//						}
//
//					case conf.ReturnBreak:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						if doRt == true {
//							isBreak = true
//							isContinue = true
//						}
//					case conf.ReturnArrStr:
//						doRt, ok := doV.([]string)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					default:
//						return rt, errors.New("each第" + strconv.Itoa(i) + "参数:不支持的数据返回类型")
//					}
//				}
//
//			}
//		}
//	case conf.ReturnInt64:
//		arrArr, _ := arr.([]int64)
//		if len(arrArr) > 0 {
//			//执行结果
//			isBreak := false
//			for i := 0; i < len(arrArr) && isBreak == false; i++ {
//				a.localParams[conf.SYS_E] = int(arrArr[i])
//				a.localParams[conf.SYS_I] = i
//				isContinue := false
//				for j := 1; j < len(doExpr) && isContinue == false; j++ {
//					a.localParams[conf.SYS_V] = rt[j]
//					doV, err := a.ExprASTResult(doExpr[j]["do"].(types.ExprAST))
//					if err != nil {
//						return nil, err
//					}
//					rType := doExpr[j]["rType"].(string)
//					switch rType {
//					case conf.ReturnInt:
//						doRt, ok := convert.Int(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnFloat:
//						doRt, ok := convert.Float64(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnStr:
//						doRt, ok := convert.String(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnBool:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnDate:
//						doRt, ok := convert.Date(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnDecimal:
//						doRt, ok := convert.Decimal(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnAny:
//						rt[i] = doV
//					case conf.ReturnMap:
//						doRt, ok := doV.(map[string]interface{})
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//
//					case conf.ReturnArray:
//						doRt, ok := eUtils.ArrType(doV, 1)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnContinue:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						if doRt == true {
//							isContinue = true
//						}
//
//					case conf.ReturnBreak:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						if doRt == true {
//							isBreak = true
//							isContinue = true
//						}
//					case conf.ReturnArrStr:
//						doRt, ok := doV.([]string)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					default:
//						return rt, errors.New("each第" + strconv.Itoa(i) + "参数:不支持的数据返回类型")
//					}
//				}
//
//			}
//		}
//	case conf.ReturnUInt8:
//		arrArr, _ := arr.([]uint8)
//		if len(arrArr) > 0 {
//			//执行结果
//			isBreak := false
//			for i := 0; i < len(arrArr) && isBreak == false; i++ {
//				a.localParams[conf.SYS_E] = int(arrArr[i])
//				a.localParams[conf.SYS_I] = i
//				isContinue := false
//				for j := 1; j < len(doExpr) && isContinue == false; j++ {
//					a.localParams[conf.SYS_V] = rt[j]
//					doV, err := a.ExprASTResult(doExpr[j]["do"].(types.ExprAST))
//					if err != nil {
//						return nil, err
//					}
//					rType := doExpr[j]["rType"].(string)
//					switch rType {
//					case conf.ReturnInt:
//						doRt, ok := convert.Int(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnFloat:
//						doRt, ok := convert.Float64(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnStr:
//						doRt, ok := convert.String(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnBool:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnDate:
//						doRt, ok := convert.Date(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnDecimal:
//						doRt, ok := convert.Decimal(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnAny:
//						rt[i] = doV
//					case conf.ReturnMap:
//						doRt, ok := doV.(map[string]interface{})
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//
//					case conf.ReturnArray:
//						doRt, ok := eUtils.ArrType(doV, 1)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnContinue:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						if doRt == true {
//							isContinue = true
//						}
//
//					case conf.ReturnBreak:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						if doRt == true {
//							isBreak = true
//							isContinue = true
//						}
//					case conf.ReturnArrStr:
//						doRt, ok := doV.([]string)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					default:
//						return rt, errors.New("each第" + strconv.Itoa(i) + "参数:不支持的数据返回类型")
//					}
//				}
//
//			}
//		}
//	case conf.ReturnUInt16:
//		arrArr, _ := arr.([]uint16)
//		if len(arrArr) > 0 {
//			//执行结果
//			isBreak := false
//			for i := 0; i < len(arrArr) && isBreak == false; i++ {
//				a.localParams[conf.SYS_E] = int(arrArr[i])
//				a.localParams[conf.SYS_I] = i
//				isContinue := false
//				for j := 1; j < len(doExpr) && isContinue == false; j++ {
//					a.localParams[conf.SYS_V] = rt[j]
//					doV, err := a.ExprASTResult(doExpr[j]["do"].(types.ExprAST))
//					if err != nil {
//						return nil, err
//					}
//					rType := doExpr[j]["rType"].(string)
//					switch rType {
//					case conf.ReturnInt:
//						doRt, ok := convert.Int(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnFloat:
//						doRt, ok := convert.Float64(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnStr:
//						doRt, ok := convert.String(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnBool:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnDate:
//						doRt, ok := convert.Date(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnDecimal:
//						doRt, ok := convert.Decimal(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnAny:
//						rt[i] = doV
//					case conf.ReturnMap:
//						doRt, ok := doV.(map[string]interface{})
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//
//					case conf.ReturnArray:
//						doRt, ok := eUtils.ArrType(doV, 1)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnContinue:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						if doRt == true {
//							isContinue = true
//						}
//
//					case conf.ReturnBreak:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						if doRt == true {
//							isBreak = true
//							isContinue = true
//						}
//					case conf.ReturnArrStr:
//						doRt, ok := doV.([]string)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					default:
//						return rt, errors.New("each第" + strconv.Itoa(i) + "参数:不支持的数据返回类型")
//					}
//				}
//
//			}
//		}
//	case conf.ReturnUInt32:
//		arrArr, _ := arr.([]uint32)
//		if len(arrArr) > 0 {
//			//执行结果
//			isBreak := false
//			for i := 0; i < len(arrArr) && isBreak == false; i++ {
//				a.localParams[conf.SYS_E] = int(arrArr[i])
//				a.localParams[conf.SYS_I] = i
//				isContinue := false
//				for j := 1; j < len(doExpr) && isContinue == false; j++ {
//					a.localParams[conf.SYS_V] = rt[j]
//					doV, err := a.ExprASTResult(doExpr[j]["do"].(types.ExprAST))
//					if err != nil {
//						return nil, err
//					}
//					rType := doExpr[j]["rType"].(string)
//					switch rType {
//					case conf.ReturnInt:
//						doRt, ok := convert.Int(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnFloat:
//						doRt, ok := convert.Float64(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnStr:
//						doRt, ok := convert.String(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnBool:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnDate:
//						doRt, ok := convert.Date(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnDecimal:
//						doRt, ok := convert.Decimal(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnAny:
//						rt[i] = doV
//					case conf.ReturnMap:
//						doRt, ok := doV.(map[string]interface{})
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//
//					case conf.ReturnArray:
//						doRt, ok := eUtils.ArrType(doV, 1)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnContinue:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						if doRt == true {
//							isContinue = true
//						}
//
//					case conf.ReturnBreak:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						if doRt == true {
//							isBreak = true
//							isContinue = true
//						}
//					case conf.ReturnArrStr:
//						doRt, ok := doV.([]string)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					default:
//						return rt, errors.New("each第" + strconv.Itoa(i) + "参数:不支持的数据返回类型")
//					}
//				}
//
//			}
//		}
//	case conf.ReturnUInt64:
//		arrArr, _ := arr.([]uint64)
//		if len(arrArr) > 0 {
//			//执行结果
//			isBreak := false
//			for i := 0; i < len(arrArr) && isBreak == false; i++ {
//				a.localParams[conf.SYS_E] = int(arrArr[i])
//				a.localParams[conf.SYS_I] = i
//				isContinue := false
//				for j := 1; j < len(doExpr) && isContinue == false; j++ {
//					a.localParams[conf.SYS_V] = rt[j]
//					doV, err := a.ExprASTResult(doExpr[j]["do"].(types.ExprAST))
//					if err != nil {
//						return nil, err
//					}
//					rType := doExpr[j]["rType"].(string)
//					switch rType {
//					case conf.ReturnInt:
//						doRt, ok := convert.Int(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnFloat:
//						doRt, ok := convert.Float64(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnStr:
//						doRt, ok := convert.String(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnBool:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnDate:
//						doRt, ok := convert.Date(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnDecimal:
//						doRt, ok := convert.Decimal(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnAny:
//						rt[i] = doV
//					case conf.ReturnMap:
//						doRt, ok := doV.(map[string]interface{})
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//
//					case conf.ReturnArray:
//						doRt, ok := eUtils.ArrType(doV, 1)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnContinue:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						if doRt == true {
//							isContinue = true
//						}
//
//					case conf.ReturnBreak:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						if doRt == true {
//							isBreak = true
//							isContinue = true
//						}
//					case conf.ReturnArrStr:
//						doRt, ok := doV.([]string)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					default:
//						return rt, errors.New("each第" + strconv.Itoa(i) + "参数:不支持的数据返回类型")
//					}
//				}
//
//			}
//		}
//	case conf.ReturnUInt:
//		arrArr, _ := arr.([]uint64)
//		if len(arrArr) > 0 {
//			//执行结果
//			isBreak := false
//			for i := 0; i < len(arrArr) && isBreak == false; i++ {
//				a.localParams[conf.SYS_E] = int(arrArr[i])
//				a.localParams[conf.SYS_I] = i
//				isContinue := false
//				for j := 1; j < len(doExpr) && isContinue == false; j++ {
//					a.localParams[conf.SYS_V] = rt[j]
//					doV, err := a.ExprASTResult(doExpr[j]["do"].(types.ExprAST))
//					if err != nil {
//						return nil, err
//					}
//					rType := doExpr[j]["rType"].(string)
//					switch rType {
//					case conf.ReturnInt:
//						doRt, ok := convert.Int(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnFloat:
//						doRt, ok := convert.Float64(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnStr:
//						doRt, ok := convert.String(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnBool:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnDate:
//						doRt, ok := convert.Date(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnDecimal:
//						doRt, ok := convert.Decimal(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnAny:
//						rt[i] = doV
//					case conf.ReturnMap:
//						doRt, ok := doV.(map[string]interface{})
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//
//					case conf.ReturnArray:
//						doRt, ok := eUtils.ArrType(doV, 1)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					case conf.ReturnContinue:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						if doRt == true {
//							isContinue = true
//						}
//
//					case conf.ReturnBreak:
//						doRt, ok := convert.Bool(doV)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						if doRt == true {
//							isBreak = true
//							isContinue = true
//						}
//					case conf.ReturnArrStr:
//						doRt, ok := doV.([]string)
//						if ok == false {
//							return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变成返回值配置!" + t)
//						}
//						rt[j] = doRt
//					default:
//						return rt, errors.New("each第" + strconv.Itoa(i) + "参数:不支持的数据返回类型")
//					}
//				}
//
//			}
//		}
//	default:
//		return nil, errors.New("不合法的数组类型!")
//	}
//	for i := 1; i < len(exprs); i++ {
//		result[keys[i]] = rt[i]
//	}
//	return result, err
//}

/*
*
foreach
*/
func (a *AST) defForeach(exprs ...types.ExprAST) (interface{}, error) {
	//判断参数
	if len(exprs) < 2 {
		return nil, errors.New("defForeach 参数数目错误")
	}
	arr, err := a.ExprASTResult(exprs[0])
	if err != nil {
		return nil, err
	}
	//判断数组类型
	t, bl := eUtils.ArrType(arr, 1)
	if bl == false {
		return nil, errors.New("defForeach 错误的参数类型")
	}

	//数据转换
	toVal, ok := eUtils.ToTypeVal(arr, t)
	if !ok {
		return nil, errors.New("不合法的数组类型!")
	}
	//判断数组类别
	switch t {
	case consts.ARR_MAP:
		return forEach(toVal.([]map[string]interface{}), a, exprs...)
	case consts.ARR_STR:
		return forEach(toVal.([]string), a, exprs...)
	case consts.ARR_ANY:
		return forEach(toVal.([]interface{}), a, exprs...)
	case consts.ARR_DECIMAL:
		return forEach(toVal.([]decimal.Decimal), a, exprs...)
	case consts.ARR_INT,
		consts.ARR_INT8,
		consts.ARR_INT16,
		consts.ARR_INT32,
		consts.ARR_INT64,
		consts.ARR_UINT8,
		consts.ARR_UINT16,
		consts.ARR_UINT32,
		consts.ARR_UINT64,
		consts.ARR_UINT:
		return forEach(toVal.([]int), a, exprs...)
	case consts.ARR_FLOAT,
		consts.ARR_FLOAT32:
		return forEach(toVal.([]float64), a, exprs...)
	case consts.ARR_BOOL:
		return forEach(toVal.([]bool), a, exprs...)
	case consts.ARR_DATE:
		return forEach(toVal.([]time.Time), a, exprs...)
	}
	return nil, errors.New("不合法的数组类型!")
}

/*
*
defEach
*/
func (a *AST) defEach(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	//转化为字符串
	str, bl := convert.String(v)
	if bl == false {
		return nil, errors.New("each函数参数1返回类型不正确，应为str!")
	}
	v2, err := a.ExprASTResult(expr[3])
	if err != nil {
		return nil, err
	}
	//转化为字符串
	str2, bl := convert.String(v2)
	if bl == false {
		return nil, errors.New("each函数参数4返回类型不正确，应为str!")
	}

	return map[string]interface{}{
		"rType": str,
		"init":  expr[1],
		"do":    expr[2],
		"key":   str2,
	}, nil
}

/*
*
泛型 foreach
*/
func forEach[Data any](arr []Data, a *AST, exprs ...types.ExprAST) (interface{}, error) {
	//声明返回值和执行表达式
	rt := make([]interface{}, len(exprs))
	keys := make(map[int]string)
	doExpr := make([]map[string]interface{}, len(exprs))
	result := make(map[string]interface{})

	//哨兵机制 避免每次循环都去判断
	for i := 1; i < len(exprs); i++ {
		v, err := a.ExprASTResult(exprs[i])
		if err != nil {
			return result, nil
		}

		res := v.(map[string]interface{})
		keys[i] = res["key"].(string)
		rType := res["rType"].(string)

		//处理init
		doV, err := a.ExprASTResult(res["init"].(types.ExprAST))
		if err != nil {
			return rt, err
		}
		toVal, ok := eUtils.ToTypeVal(doV, rType)
		if !ok {
			return rt, errors.New("foreach 不支持的返回类型:t=" + rType)
		}
		rt[i] = toVal

		//写入到任务
		doExpr[i] = res
	}

	//执行结果
	isBreak := false
	for i := 0; i < len(arr) && !isBreak; i++ {
		//写入局部变量
		a.localParams[conf.SYS_E] = arr[i]
		a.localParams[conf.SYS_I] = i
		isContinue := false
		for j := 1; j < len(doExpr) && !isContinue; j++ {
			a.localParams[conf.SYS_V] = rt[j]
			rType := doExpr[j]["rType"].(string)

			//执行do
			doV, err := a.ExprASTResult(doExpr[j]["do"].(types.ExprAST))
			if err != nil {
				return nil, err
			}
			toV, ok := eUtils.ToTypeVal(doV, rType)

			if !ok {
				return rt, errors.New("each第" + strconv.Itoa(j) + "参数 计算结果无法转变" + rType)
			}

			//是否continue
			if rType == consts.CONTINUE && toV.(string) == "continue" {
				isContinue = true
				continue
			}

			//是否break
			if rType == consts.BREAK && toV.(string) == "break" {
				isBreak = true
				isContinue = true
				continue
			}

			//继续执行
			rt[j] = toV
		}
	}
	for i := 1; i < len(exprs); i++ {
		result[keys[i]] = rt[i]
	}
	return result, nil
}
