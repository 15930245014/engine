package dsl

import (
	"errors"
	"gitlab.shudieds.com/mec/lib/utils/convert"
	"gitlab.shudieds.com/zxh/engine/types"
	"strings"
)

/*
*

	获取字符串长度
*/
func (a *AST) strLen(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	str, bl := convert.String(v)
	if bl == false {
		return nil, errors.New("strLen 参数错误!")
	}
	return len([]rune(str)), nil
}

func (a *AST) strStr(expr ...types.ExprAST) (interface{}, error) {
	return expr[0].GetVal(), nil
}

/*
*

	连接字符串
*/
func (a *AST) strConcat(expr ...types.ExprAST) (interface{}, error) {
	var result string
	if len(expr) == 0 {
		return result, errors.New("strConcat 参数不能为空!")
	}

	for i := 0; i < len(expr); i++ {
		v, err := a.ExprASTResult(expr[i])
		if err != nil {
			return "", err
		}

		str, bl := convert.String(v)
		if bl == false {
			return nil, errors.New("strConcat 参数错误")
		}
		result = result + str
	}
	return result, nil
}

/*
*

	是否在目标字符数组
*/
func (a *AST) strIn(expr ...types.ExprAST) (interface{}, error) {
	if len(expr) < 2 {
		return nil, errors.New("strIn 参数数目错误!")
	}
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	str, bl := convert.String(v1)
	if bl == false {
		return nil, errors.New("strIn 参数内容错误!")
	}

	for i := 1; i < len(expr); i++ {
		v2, err := a.ExprASTResult(expr[i])
		if err != nil {
			return nil, err
		}

		str2, bl := convert.String(v2)
		if bl == false {
			return nil, errors.New("strIn 参数错误")
		}
		if str == str2 {
			return true, nil
		}
	}
	return false, nil
}

/*
*

	是否不在目标字符数组
*/
func (a *AST) strNotIn(expr ...types.ExprAST) (interface{}, error) {
	if len(expr) < 2 {
		return nil, errors.New("strNotIn 参数数目错误!")
	}
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	str, bl := convert.String(v1)
	if bl == false {
		return nil, errors.New("strNotIn 参数内容错误!")
	}

	for i := 1; i < len(expr); i++ {
		v2, err := a.ExprASTResult(expr[i])
		if err != nil {
			return nil, err
		}

		str2, bl := convert.String(v2)
		if bl == false {
			return nil, errors.New("strNotIn 参数错误")
		}
		if str == str2 {
			return false, nil
		}
	}
	return true, nil
}

/*
*

	字符串转为大写
*/
func (a *AST) strToUpper(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	str, bl := convert.String(v)
	if bl == false {
		return nil, errors.New("strToUpper 参数错误")
	}
	return strings.ToUpper(str), nil
}

/*
*

	字符串转为小写
*/
func (a *AST) strToLower(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	str, bl := convert.String(v)
	if bl == false {
		return nil, errors.New("strToLower 参数错误")
	}
	return strings.ToLower(str), nil
}

/*
*
字符串1是否以字符串2开头
*/
func (a *AST) strHasPrefix(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}

	str1, bl1 := convert.String(v1)
	str2, bl2 := convert.String(v2)

	if bl1 == false || bl2 == false {
		return nil, errors.New("strHasPrefix 参数错误")
	}

	return strings.HasPrefix(str1, str2), nil
}

/*
*

	字符串1是否以字符串2结尾
*/
func (a *AST) strHasSuffix(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}

	str1, bl1 := convert.String(v1)
	str2, bl2 := convert.String(v2)

	if bl1 == false || bl2 == false {
		return nil, errors.New("strHasSuffix 参数错误")
	}

	return strings.HasSuffix(str1, str2), nil
}

/*
*

	字符串1是否包含字符串2
*/
func (a *AST) strContains(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}

	str1, bl1 := convert.String(v1)
	str2, bl2 := convert.String(v2)

	if bl1 == false || bl2 == false {
		return nil, errors.New("strContains 参数错误")
	}

	return strings.Contains(str1, str2), nil
}

/*
*

	字符串1是否包含字符串2
*/
func (a *AST) strCount(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}

	str1, bl1 := convert.String(v1)
	str2, bl2 := convert.String(v2)

	if bl1 == false || bl2 == false {
		return nil, errors.New("strCount 参数错误")
	}

	return strings.Count(str1, str2), nil
}

/*
*

	字符串2在字符串1第一次出现位置
*/
func (a *AST) strIndex(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}
	str1, bl1 := convert.String(v1)
	str2, bl2 := convert.String(v2)

	if bl1 == false || bl2 == false {
		return nil, errors.New("strIndex 参数错误")
	}
	// 子串在字符串的字节位置
	index := strings.Index(str1, str2)
	if index >= 0 {
		// 获得子串之前的字符串并转换成[]byte
		prefix := convert.StringToBytes(str1)[0:index]
		// 将子串之前的字符串转换成[]rune
		rs := []rune(convert.BytesToString(prefix))
		// 获得子串之前的字符串的长度，便是子串在字符串的字符位置
		index = len(rs)
	}
	return index, nil
}

/*
*

	字符串2在字符串1第一次出现位置
*/
func (a *AST) strLastIndex(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}

	str1, bl1 := convert.String(v1)
	str2, bl2 := convert.String(v2)

	if bl1 == false || bl2 == false {
		return nil, errors.New("strLastIndex 参数错误")
	}

	index := strings.LastIndex(str1, str2)
	if index >= 0 {
		// 获得子串之前的字符串并转换成[]byte
		prefix := convert.StringToBytes(str1)[0:index]
		// 将子串之前的字符串转换成[]rune
		rs := []rune(convert.BytesToString(prefix))
		// 获得子串之前的字符串的长度，便是子串在字符串的字符位置
		index = len(rs)
	}
	return index, nil
}

/*
*

	去除字符串2边空格
*/
func (a *AST) strTrimSpace(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	str, bl := convert.String(v1)
	if bl == false {
		return "", errors.New("strTrimSpace 参数错误")
	}

	return strings.TrimSpace(str), nil
}

/*
*

	去除str1两边的str2
*/
func (a *AST) strTrim(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}
	str1, bl1 := convert.String(v1)
	str2, bl2 := convert.String(v2)

	if bl1 == false {
		return "", errors.New("strTrim 参数错误")
	}
	if bl2 == false {
		return nil, errors.New("strTrim 参数错误")
	}
	return strings.Trim(str1, str2), nil
}

/*
*

	去除str1左边的多个str2
*/
func (a *AST) strTrimLeft(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}

	str1, bl1 := convert.String(v1)
	str2, bl2 := convert.String(v2)

	if bl1 == false || bl2 == false {
		return nil, errors.New("strTrimLeft 参数错误")
	}
	return strings.TrimLeft(str1, str2), nil
}

/*
*

	去除str1右边的多个str2
*/
func (a *AST) strTrimRight(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}
	str1, bl1 := convert.String(v1)
	str2, bl2 := convert.String(v2)

	if bl1 == false || bl2 == false {
		return nil, errors.New("strTrimRight 参数错误")
	}
	return strings.TrimRight(str1, str2), nil
}

/*
*

	指定前缀trim
*/
func (a *AST) strTrimPrefix(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}
	str1, bl1 := convert.String(v1)
	str2, bl2 := convert.String(v2)

	if bl1 == false || bl2 == false {
		return nil, errors.New("strTrimPrefix 参数错误")
	}
	return strings.TrimPrefix(str1, str2), nil
}

/*
*

	指定后缀trim
*/
func (a *AST) strTrimSuffix(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}
	str1, bl1 := convert.String(v1)
	str2, bl2 := convert.String(v2)

	if bl1 == false || bl2 == false {
		return nil, errors.New("strTrimSuffix 参数错误")
	}
	return strings.TrimSuffix(str1, str2), nil
}

/*
*

	字符串替换
*/
func (a *AST) strReplace(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}
	v3, err := a.ExprASTResult(expr[2])
	if err != nil {
		return nil, err
	}
	str1, bl1 := convert.String(v1)
	str2, bl2 := convert.String(v2)
	str3, bl3 := convert.String(v3)

	if bl1 == false || bl2 == false || bl3 == false {
		return nil, errors.New("strReplace 参数错误")
	}
	return strings.Replace(str1, str2, str3, -1), nil
}

/*
*

	split
*/
func (a *AST) strSplit(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}
	str1, bl1 := convert.String(v1)
	str2, bl2 := convert.String(v2)

	if bl1 == false || bl2 == false {
		return nil, errors.New("strSplit 参数错误!")
	}
	return strings.Split(str1, str2), nil
}

/*
*

	join
*/
func (a *AST) strJoin(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	arr, bl := v.([]string)
	if bl == false {
		return "", errors.New("strJoin 参数错误!")
	}
	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}
	str2, _ := convert.String(v2)
	return strings.Join(arr, str2), nil
}

/*
*

	字符串重复次数
*/
func (a *AST) strRepeat(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}
	str, bl1 := convert.String(v1)
	num, bl2 := convert.Int(v2)
	if bl1 == false || bl2 == false {
		return nil, errors.New("strRepeat 参数错误!")
	}
	return strings.Repeat(str, num), nil
}

/*
*

	字符串比较
*/
func (a *AST) strCompare(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}
	str1, bl1 := convert.String(v1)
	str2, bl2 := convert.String(v2)
	if bl1 == false || bl2 == false {
		return nil, errors.New("strCompare 参数错误!")
	}
	return strings.Compare(str1, str2), err
}

func (a *AST) strConcatFill(expr ...types.ExprAST) (interface{}, error) {
	var result string
	if len(expr) == 0 {
		return result, errors.New("strConcatFill 参数不能为空!")
	}

	fistVal, err := a.ExprASTResult(expr[0])
	if err != nil {
		return "", err
	}
	fistStr, bl := convert.String(fistVal)
	if bl == false {
		return nil, errors.New("strConcatFill 参数错误")
	}
	for i := 1; i < len(expr); i++ {
		v, err := a.ExprASTResult(expr[i])
		if err != nil {
			return "", err
		}
		if v == nil {
			result = result + fistStr
			continue
		}
		str, _ := convert.String(v)
		if str == "" {
			str = fistStr
		}
		result = result + str
	}
	return result, nil
}

func (a *AST) strSub(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, errors.New("strSub 参数1计算错误" + err.Error())
	}
	str, bl := convert.String(v1)
	if bl == false {
		return nil, errors.New("strSub 参数1无法转换成str!")
	}

	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, errors.New("strSub 参数2计算错误" + err.Error())
	}
	s, bl := convert.Int(v2)
	if bl == false {
		return nil, errors.New("strSub 参数2无法转换成int!")
	}

	v3, err := a.ExprASTResult(expr[2])
	if err != nil {
		return nil, err
	}
	length, bl := convert.Int(v3)
	if bl == false {
		return nil, errors.New("strSub 参数3无法转换成int!")
	}
	runes := []rune(str)
	n := len(runes)
	if s < 0 {
		s = 0
	}
	e := s + length
	if length < 0 || e > n {
		e = n
	}
	//左开右闭ni
	return string(runes[s:e]), nil
}
func (a *AST) strLastCut(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}

	str1, bl1 := convert.String(v1)
	str2, bl2 := convert.String(v2)

	if bl1 == false || bl2 == false {
		return nil, errors.New("strLastCut 参数错误")
	}
	if len(str2) == 0 {
		return str1, nil
	}
	index := strings.LastIndex(str1, str2)
	if index < 0 {
		return str1, nil
	}

	// 获得子串之前的字符串并转换成[]byte
	prefix := convert.StringToBytes(str1)[0:index]
	// 将子串之前的字符串转换成[]rune
	rs := []rune(convert.BytesToString(prefix))
	// 获得子串之前的字符串的长度，便是子串在字符串的字符位置
	index = len(rs)

	//返回截取字符串
	return string([]rune(str1)[:index]), nil
}
