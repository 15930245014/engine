package utils

import (
	"gitlab.shudieds.com/mec/lib/consts"
	"gitlab.shudieds.com/mec/lib/entry/engine"
)

/**
任意类型转化
*/

func ToTypeVal(v interface{}, t string) (interface{}, bool) {
	//判断类型是否合法
	if _, ok := CVT[t]; !ok {
		return nil, false
	}

	//判断是否为树形结构
	if tr, ok := v.(*engine.TreeParams); ok {
		v = tr.Data
	}

	return CVT[t](v)
}

func ToMap(v interface{}) (map[string]interface{}, bool) {
	val, ok := ToTypeVal(v, consts.MAP)
	if !ok {
		return nil, false
	}
	return val.(map[string]interface{}), true
}
