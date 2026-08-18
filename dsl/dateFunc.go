package dsl

import (
	"errors"
	"gitlab.shudieds.com/zxh/engine/utils"
	"strconv"
	"strings"
	"time"

	"gitlab.shudieds.com/mec/lib/utils/convert"
	"gitlab.shudieds.com/zxh/engine/types"
)

func (a *AST) dateTimestamp(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	t, bl := v1.(time.Time)
	if bl == false {
		return nil, errors.New("dateTimestamp 参数错误!")
	}
	return t.Unix(), nil
}

func (a *AST) dateTimeFormat(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}

	t, bl1 := convert.Date(v1)
	f, bl2 := convert.String(v2)

	if bl1 == false || bl2 == false {
		return nil, errors.New("dateTimeFormat 参数错误!")
	}
	if f == "" {
		f = "2006-01-02 15:04:05"
	}

	return t.Format(f), nil
}

func (a *AST) dateGetYear(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	t, bl := v1.(time.Time)
	if bl == false {
		return nil, errors.New("dateGetYear 参数错误!")
	}
	return t.Format("2006"), nil
}

func (a *AST) dateGetMonth(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	t, bl := convert.Date(v1)
	if bl == false {
		return nil, errors.New("dateGetMonth 参数错误!")
	}
	return t.Format("01"), nil
}

func (a *AST) dateGetDay(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	t, bl := convert.Date(v1)
	if bl == false {
		return nil, errors.New("dateGetDay 参数错误!")
	}
	return t.Format("02"), nil
}
func (a *AST) dateGetHour(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	t, bl := convert.Date(v1)
	if bl == false {
		return nil, errors.New("dateGetHour 参数错误!")
	}
	return t.Format("15"), nil
}

func (a *AST) dateGetMinute(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	t, bl := convert.Date(v1)
	if bl == false {
		return nil, errors.New("dateGetMinute 参数错误!")
	}
	return t.Format("04"), nil
}

func (a *AST) dateGetSecond(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	t, bl := convert.Date(v1)
	if bl == false {
		return nil, errors.New("dateGetSecond 参数错误!")
	}
	return t.Format("05"), nil
}

func (a *AST) getTimestamp(expr ...types.ExprAST) (interface{}, error) {
	return time.Now().Unix(), nil
}

func (a *AST) getCurrentTime(expr ...types.ExprAST) (interface{}, error) {
	return time.Now(), nil
}

func (a *AST) getYYmm(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	t1, bl := convert.Date(v1)
	if bl == false {
		return nil, errors.New("getYYmm 参数1错误!")
	}
	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}
	t2, bl := convert.String(v2)
	if bl == false {
		return nil, errors.New("getYYmm 参数2错误!")
	}
	v3, err := a.ExprASTResult(expr[2])
	if err != nil {
		return nil, err
	}
	t3, bl := convert.Int(v3)
	if bl == false {
		return nil, errors.New("getYYmm 参数3错误!")
	}
	years := t3 / 12
	months := t3 % 12

	// get datetime parts
	ye := t1.Year()
	mo := t1.Month()

	// years
	ye += years

	// months
	mo += time.Month(months)
	if mo > 12 {
		mo -= 12
		ye++
	} else if mo < 1 {
		mo += 12
		ye--
	}
	if mo < 10 {
		return strconv.Itoa(ye) + t2 + "0" + strconv.Itoa(int(mo)), nil
	}
	return strconv.Itoa(ye) + t2 + strconv.Itoa(int(mo)), nil
}

// 获取毫秒时间戳
func (a *AST) getUnixMilliTime(expr ...types.ExprAST) (interface{}, error) {
	return time.Now().UnixMilli(), nil
}

// 获取毫秒时间戳
func (a *AST) dateGetLastMonthDay(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	t1, bl := convert.Date(v1)
	if bl == false {
		return nil, errors.New("dateGetLastMonthDay 参数错误!")
	}
	dateStr := t1.Format("2006-01")
	dateStrArr := strings.Split(dateStr, `-`)
	//fmt.Println(dateStrArr)
	_, last := utils.GetMonthStartAndEnd(dateStrArr[0], dateStrArr[1])
	return last, nil

}
