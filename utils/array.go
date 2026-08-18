package utils

import (
	"gitlab.shudieds.com/mec/lib/utils/convert"
)

/**
获取数组类别
*/

func ArrType(v interface{}, deep int) (string, bool) {
	return convert.ArrType(v, deep)

}

func ToArrInterface(v interface{}) ([]interface{}, bool) {
	return convert.ToArrAny(v)
}
