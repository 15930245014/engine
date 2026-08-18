package utils

import (
	"errors"
	"gitlab.shudieds.com/mec/lib/utils/convert"
	"sort"
)

func ArrLen[Data any](datas []Data) (interface{}, error) {
	return len(datas), nil
}

func ArrIndex[Data any](arr []Data, i int) (interface{}, error) {
	if len(arr) == 0 {
		return nil, errors.New("arrIndex 数组为空!")
	}
	if i < 0 {
		i = 0
	}
	if i >= len(arr) {
		i = len(arr) - 1
	}
	return arr[i], nil
}

func ArrFirst[Data any](arr []Data) (interface{}, error) {
	if len(arr) == 0 {
		return nil, errors.New("arrIndex 数组为空!")
	}
	return arr[0], nil
}

func ArrLast[Data any](arr []Data) (interface{}, error) {
	if len(arr) == 0 {
		return nil, errors.New("arrIndex 数组为空!")
	}
	return arr[len(arr)-1], nil
}

func ArrSlice[Data any](arr []Data, s, e int) (interface{}, error) {
	if len(arr) == 0 {
		return nil, errors.New("ArrSlice 数组为空!")
	}
	if s < 0 {
		s = 0
	}
	if e > len(arr) {
		e = len(arr)
	}
	if s > e {
		return nil, errors.New("arrSlice 参数错误! 截取start > end!")
	}
	return arr[s:e], nil
}

func ArrAppend[Data any](arr []Data, val interface{}, t string) (interface{}, error) {
	toVal, ok := ToTypeVal(val, t)
	if !ok {
		return nil, errors.New("arrAppend 写入数据类型不匹配t=" + t)
	}
	arr = append(arr, toVal.(Data))
	return arr, nil
}

func ArrUnshift[Data any](arr []Data, val interface{}, t string) (interface{}, error) {
	toVal, ok := ToTypeVal(val, t)
	if !ok {
		return nil, errors.New("arrUnshift 写入数据类型不匹配t=" + t)
	}
	return append([]Data{toVal.(Data)}, arr...), nil
}

func ArrMerge[Data any](arr1, arr2 []Data) (interface{}, error) {
	return append(arr1, arr2...), nil
}

func ArrSearch[Data any](arr []Data, val any, t string) (interface{}, error) {
	toVal, ok := ToTypeVal(val, t)
	if !ok {
		return -1, errors.New("ArrSearch 目标值和数组类别不合法")
	}
	for i := 0; i < len(arr); i++ {
		res, err := convert.Compare(toVal, arr[i])
		if err == nil && res == 0 {
			return i, nil
		}
	}
	return -1, nil
}

func ArrUnique[Data any](arr []Data) (interface{}, error) {
	set := make(map[string]bool)
	rt := []Data{}
	for i := 0; i < len(arr); i++ {
		key, ok := convert.String(arr[i])
		if !ok {
			continue
		}
		if !set[key] {
			rt = append(rt, arr[i])
		}
		set[key] = true
	}
	return rt, nil
}

func ArrReverse[Data any](arr []Data) (interface{}, error) {
	rt := make([]Data, len(arr))
	for i := len(arr) - 1; i >= 0; i-- {
		rt[i] = arr[i]
	}
	return rt, nil
}

func ArrSortAsc[Data any](arr []Data) (interface{}, error) {
	if len(arr) == 0 {
		return arr, nil
	}
	sort.Slice(arr, func(i, j int) bool {
		res, _ := convert.Compare(arr[i], arr[j])
		return res == -1
	})
	return arr, nil
}

func ArrSortDesc[Data any](arr []Data) (interface{}, error) {
	if len(arr) == 0 {
		return arr, nil
	}
	sort.Slice(arr, func(i, j int) bool {
		res, _ := convert.Compare(arr[i], arr[j])
		return res == 1
	})
	return arr, nil
}
