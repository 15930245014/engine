package dsl

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.shudieds.com/mec/lib/consts"
	eUtils "gitlab.shudieds.com/zxh/engine/utils"
	"sort"
	"strings"

	"gitlab.shudieds.com/mec/lib/utils/convert"
	"gitlab.shudieds.com/mec/lib/utils/maphelper"
	"gitlab.shudieds.com/zxh/engine/types"
)

func (a *AST) mapLen(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	mp, bl := eUtils.ToMap(v1)
	if bl == false {
		return nil, errors.New("mapLen 参数错误!")
	}
	return len(mp), nil
}

func (a *AST) mapIsEmpty(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	mp, bl := eUtils.ToMap(v1)
	if bl == false {
		return true, nil
	}
	return len(mp) == 0, nil
}

func (a *AST) mapKeyExist(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}

	mp, bl1 := eUtils.ToMap(v1)
	key, bl2 := convert.String(v2)
	if bl1 == false || bl2 == false {
		return nil, errors.New("mapKeyExist 参数错误!")
	}

	_, ok := mp[key]
	return ok, nil
}

func (a *AST) mapPut(expr ...types.ExprAST) (interface{}, error) {
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
	mp, bl1 := eUtils.ToMap(v1)
	k, bl2 := convert.String(v2)
	if bl1 == false || bl2 == false {
		return nil, errors.New("mapPut 参数错误!")
	}
	mp[k] = v3
	return mp, nil
}

func (a *AST) mapGet(expr ...types.ExprAST) (interface{}, error) {
	var v3 interface{}
	//获取默认值
	if len(expr) >= 3 {
		v3, _ = a.ExprASTResult(expr[2])
	}

	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		if v3 != nil {
			return v3, nil
		}
		return nil, err
	}
	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		if v3 != nil {
			return v3, nil
		}
		return nil, err
	}

	mp, bl1 := eUtils.ToMap(v1)
	if !bl1 {
		if v3 != nil {
			return v3, nil
		}
		return nil, errors.New("mapGet 参数错误!")
	}

	k, bl2 := convert.String(v2)
	if bl2 {
		val, bl := maphelper.JDGet(mp, []string{k})
		if bl == false {
			if v3 != nil {
				return v3, nil
			}
			return nil, errors.New("key不存在")
		}
		return val, nil
	}
	kArr, bl2 := eUtils.ToTypeVal(v2, consts.ARR_STR)
	if !bl2 {
		return nil, errors.New("mapGet 参数错误 key为str或arr.str")
	}
	val, bl := maphelper.JDGet(mp, kArr.([]string))
	if bl == false {
		if v3 != nil {
			return v3, nil
		}
		return nil, errors.New("key值不存在")
	}
	return val, nil

}

func (a *AST) mapKeys(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	mp, bl := eUtils.ToMap(v1)
	if bl == false {
		return nil, errors.New("mapKeys 参数错误!")
	}
	if len(mp) == 0 {
		return []string{}, nil
	}
	var rt []string
	for k, _ := range mp {
		rt = append(rt, k)
	}
	return rt, nil
}

func (a *AST) mapValues(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	mp, bl := eUtils.ToMap(v1)
	if bl == false {
		return nil, errors.New("mapValues 参数错误!")
	}

	if bl == false {
		return nil, errors.New("mapKeys 参数错误!")
	}

	var rt []interface{}
	if len(mp) > 0 {
		for _, v := range mp {
			rt = append(rt, v)
		}
	}
	return rt, nil
}

func (a *AST) mapMerge(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}

	mp1, bl := eUtils.ToMap(v1)
	mp2, b2 := eUtils.ToMap(v2)
	if bl == false || b2 == false {
		return nil, errors.New("mapMerge 参数错误!")
	}
	if len(mp1) == 0 {
		return mp2, nil
	}
	if len(mp2) == 0 {
		return mp1, nil
	}

	for k, v := range mp2 {
		if _, ok := mp1[k]; !ok {
			mp1[k] = v
		}
	}
	return mp1, nil
}

func (a *AST) mapNestedGet(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}

	mp, bl1 := eUtils.ToMap(v1)
	k, bl2 := convert.String(v2)
	if bl1 == false || bl2 == false {
		return nil, errors.New("mapNestedGet 参数错误!")
	}
	kArr := strings.Split(k, ".")

	val, bl := maphelper.JDGet(mp, kArr)
	if bl == false {
		return nil, nil
	}
	return val, nil
}

func (a *AST) mapSet(expr ...types.ExprAST) (interface{}, error) {
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
	mp, bl1 := eUtils.ToMap(v1)
	k, bl2 := convert.String(v2)
	if bl1 == false || bl2 == false {
		return nil, errors.New("mapNestedSet 参数错误!")
	}
	kArr := strings.Split(k, ".")

	ok := maphelper.JDSet(mp, kArr, v3)
	if !ok {
		return nil, errors.New("mapSet 错误!")
	}
	return mp, nil
}

func (a *AST) mapKeySortAsc(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	mp, bl := eUtils.ToMap(v1)
	if bl == false {
		return nil, errors.New("mapKeySortAsc 参数错误")
	}
	if len(mp) <= 1 {
		return mp, nil
	}
	//获取keys
	var keys []string
	for key, _ := range mp {
		keys = append(keys, key)
	}

	//key排序
	sort.Strings(keys)

	//声明返回值
	result := make(map[string]interface{})
	for _, k := range keys {
		result[k] = mp[k]
	}
	return result, nil
}

func (a *AST) mapKeySortDesc(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	mp, bl := eUtils.ToMap(v1)
	if bl == false {
		return nil, errors.New("mapKeySortDesc 参数错误")
	}
	if len(mp) <= 1 {
		return mp, nil
	}
	//获取keys
	var keys []string
	for key, _ := range mp {
		keys = append(keys, key)
	}

	//key排序
	sort.Strings(keys)

	//声明返回值
	result := make(map[string]interface{})
	for i := len(keys) - 1; i >= 0; i-- {
		result[keys[i]] = mp[keys[i]]
	}
	return result, nil
}

func (a *AST) mapRef(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	_, bl := eUtils.ToMap(v)
	if bl == false {
		return nil, errors.New("mapRef 参数值不为map类型")
	}

	return v, nil
}

// 将map 转为 字符串
// 输入
//
//	data := map[string]interface{}{
//		"page_size":             2,
//		"page":                  1,
//		"get_detail":            1,
//		"get_address":           1,
//		"get_custom_order_type": 1,
//		"condition": map[string]interface{}{
//			"order_code_list": []string{"CU658020081030002", "183115093750-1956032237008"},
//		},
//	}
//
// 输出 "{\"page_size\":2,\"page\":1,\"get_detail\":1,\"get_address\":1,\"get_custom_order_type\":1,\"condition\":{\"order_code_list\":[\"CU658020081030002\",\"183115093750-1956032237008\"]}}"
// 测试 mapToStr($ENV.ECangBizContent)
func (a *AST) mapToStr(expr ...types.ExprAST) (interface{}, error) {

	mapParams, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	data, ok := eUtils.ToMap(mapParams)
	if !ok {
		return nil, errors.New("mapToStr mapParams is 不合法")
	}
	// 将 map 转换为 JSON 字符串
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("mapToStr json Marshal 解析参数错误:%v", err)
	}

	// 将 JSON 字符串转换为包含转义字符的字符串
	escapedJSONData := fmt.Sprintf("%q", string(jsonData))

	return escapedJSONData, nil
}
