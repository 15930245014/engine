package utils

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"gitlab.shudieds.com/mec/lib/consts"
	"gitlab.shudieds.com/mec/lib/utils/convert"
	"sync"
	"time"
)

var CVT map[string]func(interface{}) (interface{}, bool)

func initCvt() {
	CVT = map[string]func(interface{}) (interface{}, bool){
		consts.INT: func(i interface{}) (interface{}, bool) {
			return convert.Int(i)
		},

		consts.BOOL: func(i interface{}) (interface{}, bool) {
			return convert.Bool(i)
		},

		consts.STR: func(i interface{}) (interface{}, bool) {
			return convert.String(i)
		},

		consts.DATE: func(i interface{}) (interface{}, bool) {
			return convert.Date(i)
		},

		consts.DECIMAL: func(i interface{}) (interface{}, bool) {
			return convert.Decimal(i)
		},

		consts.FLOAT: func(i interface{}) (interface{}, bool) {
			return convert.Float64(i)
		},

		consts.ANY: func(i interface{}) (interface{}, bool) {
			if i != nil {
				return i, true
			}
			return nil, false
		},

		consts.INT8: func(i interface{}) (interface{}, bool) {
			return convert.Int(i)
		},

		consts.INT16: func(i interface{}) (interface{}, bool) {
			return convert.Int(i)
		},

		consts.INT32: func(i interface{}) (interface{}, bool) {
			return convert.Int(i)
		},

		consts.INT64: func(i interface{}) (interface{}, bool) {
			return convert.Int(i)
		},

		consts.UINT8: func(i interface{}) (interface{}, bool) {
			return convert.Int(i)
		},

		consts.UINT16: func(i interface{}) (interface{}, bool) {
			return convert.Int(i)
		},

		consts.UINT32: func(i interface{}) (interface{}, bool) {
			return convert.Int(i)
		},

		consts.UINT64: func(i interface{}) (interface{}, bool) {
			return convert.Int(i)
		},

		consts.UINT: func(i interface{}) (interface{}, bool) {
			return convert.Int(i)
		},

		consts.FLOAT32: func(i interface{}) (interface{}, bool) {
			return convert.Float64(i)
		},

		consts.NIL: func(i interface{}) (interface{}, bool) {
			return nil, true
		},

		consts.ZERO: func(i interface{}) (interface{}, bool) {
			return "__DSL-ZERO__", true
		},

		consts.MAP: func(i interface{}) (interface{}, bool) {
			return convert.ToMap(i)
		},

		consts.ARR: func(i interface{}) (interface{}, bool) {
			return convert.ToArrAny(i)
		},

		consts.ARR_INT: func(i interface{}) (interface{}, bool) {
			arr, ok := i.([]int)
			if ok {
				return arr, true
			}
			if v, ok := i.(string); ok && v == "__DSL-ZERO__" {
				return []int{}, true
			}

			arr2, ok := i.([]interface{})
			if !ok {
				return nil, false
			}
			rt := []int{}
			for j := 0; j < len(arr2); j++ {
				toVal, ok := convert.Int(arr2[j])
				if ok {
					rt = append(rt, toVal)
				}
			}
			if len(rt) == 0 {
				return nil, false
			}

			arr3, ok := i.([]json.Number)
			if !ok {
				return nil, false
			}
			for j := 0; j < len(arr3); j++ {
				toVal, ok := convert.Int(arr2[j])
				if ok {
					rt = append(rt, toVal)
				}
			}
			if len(rt) == 0 {
				return nil, false
			}

			return rt, true
		},

		consts.ARR_STR: func(i interface{}) (interface{}, bool) {
			arr, ok := i.([]string)
			if ok {
				return arr, true
			}
			if v, ok := i.(string); ok && v == "__DSL-ZERO__" {
				return []string{}, true
			}
			arr2, ok := i.([]interface{})
			if !ok {
				return nil, false
			}
			rt := []string{}
			for j := 0; j < len(arr2); j++ {
				toVal, ok := convert.String(arr2[j])
				if ok {
					rt = append(rt, toVal)
				}
			}
			if len(rt) == 0 {
				return nil, false
			}

			arr3, ok := i.([]json.Number)
			if !ok {
				return nil, false
			}
			for j := 0; j < len(arr3); j++ {
				toVal, ok := convert.String(arr2[j])
				if ok {
					rt = append(rt, toVal)
				}
			}
			if len(rt) == 0 {
				return nil, false
			}
			return rt, true
		},

		consts.ARR_DECIMAL: func(i interface{}) (interface{}, bool) {
			//decimal
			arr, ok := i.([]decimal.Decimal)
			if ok {
				return arr, true
			}

			//字符串
			if v, ok := i.(string); ok && v == "__DSL-ZERO__" {
				return []decimal.Decimal{}, true
			}

			//any
			arr2, ok := i.([]interface{})
			if !ok {
				return nil, false
			}
			rt := []decimal.Decimal{}
			for j := 0; j < len(arr2); j++ {
				toVal, ok := convert.Decimal(arr2[j])
				if ok {
					rt = append(rt, toVal)
				}
			}
			if len(rt) == 0 {
				return nil, false
			}

			//json-number
			arr3, ok := i.([]json.Number)
			if !ok {
				return nil, false
			}
			for j := 0; j < len(arr3); j++ {
				toVal, ok := convert.Decimal(arr2[j])
				if ok {
					rt = append(rt, toVal)
				}
			}
			if len(rt) == 0 {
				return nil, false
			}

			return rt, true
		},

		consts.ARR_FLOAT: func(i interface{}) (interface{}, bool) {
			arr, ok := i.([]float64)
			if ok {
				return arr, true
			}

			if v, ok := i.(string); ok && v == "__DSL-ZERO__" {
				return []float64{}, true
			}

			arr2, ok := i.([]interface{})
			if !ok {
				return nil, false
			}
			rt := []float64{}
			for j := 0; j < len(arr2); j++ {
				toVal, ok := convert.Float64(arr2[j])
				if ok {
					rt = append(rt, toVal)
				}
			}
			if len(rt) == 0 {
				return nil, false
			}

			//json-number
			arr3, ok := i.([]json.Number)
			if !ok {
				return nil, false
			}
			for j := 0; j < len(arr3); j++ {
				toVal, ok := convert.Float64(arr2[j])
				if ok {
					rt = append(rt, toVal)
				}
			}
			if len(rt) == 0 {
				return nil, false
			}

			return rt, true
		},

		consts.ARR_DATE: func(i interface{}) (interface{}, bool) {
			arr, ok := i.([]time.Time)
			if ok {
				return arr, true
			}

			if v, ok := i.(string); ok && v == "__DSL-ZERO__" {
				return []time.Time{}, true
			}

			arr2, ok := i.([]interface{})
			if !ok {
				return nil, false
			}
			rt := []time.Time{}
			for j := 0; j < len(arr2); j++ {
				toVal, ok := convert.Date(arr2[j])
				if ok {
					rt = append(rt, toVal)
				}
			}
			if len(rt) == 0 {
				return nil, false
			}

			//json-number
			arr3, ok := i.([]json.Number)
			if !ok {
				return nil, false
			}
			for j := 0; j < len(arr3); j++ {
				toVal, ok := convert.Date(arr2[j])
				if ok {
					rt = append(rt, toVal)
				}
			}
			if len(rt) == 0 {
				return nil, false
			}

			return rt, true
		},

		consts.ARR_BOOL: func(i interface{}) (interface{}, bool) {
			arr, ok := i.([]bool)
			if ok {
				return arr, true
			}
			if v, ok := i.(string); ok && v == "__DSL-ZERO__" {
				return []bool{}, true
			}

			arr2, ok := i.([]interface{})
			if !ok {
				return nil, false
			}
			rt := []bool{}
			for j := 0; j < len(arr2); j++ {
				toVal, ok := convert.Bool(arr2[j])
				if ok {
					rt = append(rt, toVal)
				}
			}
			if len(rt) == 0 {
				return nil, false
			}

			//json-number
			arr3, ok := i.([]json.Number)
			if !ok {
				return nil, false
			}
			for j := 0; j < len(arr3); j++ {
				toVal, ok := convert.Bool(arr2[j])
				if ok {
					rt = append(rt, toVal)
				}
			}
			if len(rt) == 0 {
				return nil, false
			}

			return rt, true
		},

		consts.ARR_ANY: func(i interface{}) (interface{}, bool) {
			return convert.ToArrAny(i)
		},

		consts.ARR_MAP: func(i interface{}) (interface{}, bool) {
			return convert.ToArrMap(i)
		},

		consts.ARR_INT8: func(i interface{}) (interface{}, bool) {
			rt := []int{}
			arr, ok := i.([]int8)
			if ok {
				for j := 0; j < len(arr); j++ {
					rt = append(rt, int(arr[j]))
				}
				return rt, true
			}
			if v, ok := i.(string); ok && v == "__DSL-ZERO__" {
				return []int{}, true
			}
			arr2, ok := i.([]interface{})
			if !ok {
				return nil, false
			}

			for j := 0; j < len(arr2); j++ {
				toVal, ok := convert.Int(arr2[j])
				if ok {
					rt = append(rt, toVal)
				}
			}
			if len(rt) == 0 && len(arr2) > 0 {
				return nil, false
			}

			//json-number
			arr3, ok := i.([]json.Number)
			if !ok {
				return nil, false
			}
			for j := 0; j < len(arr3); j++ {
				toVal, ok := convert.Int(arr2[j])
				if ok {
					rt = append(rt, toVal)
				}
			}
			if len(rt) == 0 && len(arr3) > 0 {
				return nil, false
			}

			return rt, true
		},

		consts.ARR_INT16: func(i interface{}) (interface{}, bool) {
			rt := []int{}
			arr, ok := i.([]int16)
			if ok {
				for j := 0; j < len(arr); j++ {
					rt = append(rt, int(arr[j]))
				}
				return rt, true
			}
			if v, ok := i.(string); ok && v == "__DSL-ZERO__" {
				return []int{}, true
			}
			arr2, ok := i.([]interface{})
			if !ok {
				return nil, false
			}

			for j := 0; j < len(arr2); j++ {
				toVal, ok := convert.Int(arr2[j])
				if ok {
					rt = append(rt, toVal)
				}
			}
			if len(rt) == 0 && len(arr2) > 0 {
				return nil, false
			}

			//json-number
			arr3, ok := i.([]json.Number)
			if !ok {
				return nil, false
			}
			for j := 0; j < len(arr3); j++ {
				toVal, ok := convert.Int(arr2[j])
				if ok {
					rt = append(rt, toVal)
				}
			}
			if len(rt) == 0 && len(arr3) > 0 {
				return nil, false
			}

			return rt, true
		},

		consts.ARR_INT32: func(i interface{}) (interface{}, bool) {
			rt := []int{}
			arr, ok := i.([]int32)
			if ok {
				for j := 0; j < len(arr); j++ {
					rt = append(rt, int(arr[j]))
				}
				return rt, true
			}
			if v, ok := i.(string); ok && v == "__DSL-ZERO__" {
				return []int{}, true
			}
			arr2, ok := i.([]interface{})
			if !ok {
				return nil, false
			}

			for j := 0; j < len(arr2); j++ {
				toVal, ok := convert.Int(arr2[j])
				if ok {
					rt = append(rt, toVal)
				}
			}
			if len(rt) == 0 && len(arr2) > 0 {
				return nil, false
			}

			//json-number
			arr3, ok := i.([]json.Number)
			if !ok {
				return nil, false
			}
			for j := 0; j < len(arr3); j++ {
				toVal, ok := convert.Int(arr2[j])
				if ok {
					rt = append(rt, toVal)
				}
			}
			if len(rt) == 0 && len(arr3) > 0 {
				return nil, false
			}
			return rt, true
		},

		consts.ARR_INT64: func(i interface{}) (interface{}, bool) {
			rt := []int{}
			arr, ok := i.([]int64)
			if ok {
				for j := 0; j < len(arr); j++ {
					rt = append(rt, int(arr[j]))
				}
				return rt, true
			}
			if v, ok := i.(string); ok && v == "__DSL-ZERO__" {
				return []int{}, true
			}
			arr2, ok := i.([]interface{})
			if !ok {
				return nil, false
			}

			for j := 0; j < len(arr2); j++ {
				toVal, ok := convert.Int(arr2[j])
				if ok {
					rt = append(rt, toVal)
				}
			}
			if len(rt) == 0 && len(arr2) > 0 {
				return nil, false
			}

			//json-number
			arr3, ok := i.([]json.Number)
			if !ok {
				return nil, false
			}
			for j := 0; j < len(arr3); j++ {
				toVal, ok := convert.Int(arr2[j])
				if ok {
					rt = append(rt, toVal)
				}
			}
			if len(rt) == 0 && len(arr3) > 0 {
				return nil, false
			}
			return rt, true
		},

		consts.ARR_UINT: func(i interface{}) (interface{}, bool) {
			rt := []int{}
			arr, ok := i.([]uint)
			if ok {
				for j := 0; j < len(arr); j++ {
					rt = append(rt, int(arr[j]))
				}
				return rt, true
			}
			if v, ok := i.(string); ok && v == "__DSL-ZERO__" {
				return []int{}, true
			}
			arr2, ok := i.([]interface{})
			if !ok {
				return nil, false
			}

			for j := 0; j < len(arr2); j++ {
				toVal, ok := convert.Int(arr2[j])
				if ok {
					rt = append(rt, toVal)
				}
			}
			if len(rt) == 0 && len(arr2) > 0 {
				return nil, false
			}

			//json-number
			arr3, ok := i.([]json.Number)
			if !ok {
				return nil, false
			}

			for j := 0; j < len(arr3); j++ {
				toVal, ok := convert.Int(arr2[j])
				if ok {
					rt = append(rt, toVal)
				}
			}
			if len(rt) == 0 && len(arr3) > 0 {
				return nil, false
			}

			return rt, true
		},

		consts.ARR_UINT8: func(i interface{}) (interface{}, bool) {
			rt := []int{}
			arr, ok := i.([]uint8)
			if ok {
				for j := 0; j < len(arr); j++ {
					rt = append(rt, int(arr[j]))
				}
				return rt, true
			}
			if v, ok := i.(string); ok && v == "__DSL-ZERO__" {
				return []int{}, true
			}
			arr2, ok := i.([]interface{})
			if !ok {
				return nil, false
			}

			for j := 0; j < len(arr2); j++ {
				toVal, ok := convert.Int(arr2[j])
				if ok {
					rt = append(rt, toVal)
				}
			}
			if len(rt) == 0 && len(arr2) > 0 {
				return nil, false
			}

			//json-number
			arr3, ok := i.([]json.Number)
			if !ok {
				return nil, false
			}
			for j := 0; j < len(arr3); j++ {
				toVal, ok := convert.Int(arr2[j])
				if ok {
					rt = append(rt, toVal)
				}
			}
			if len(rt) == 0 && len(arr3) > 0 {
				return nil, false
			}
			return rt, true
		},

		consts.ARR_UINT16: func(i interface{}) (interface{}, bool) {
			rt := []int{}
			arr, ok := i.([]uint16)
			if ok {
				for j := 0; j < len(arr); j++ {
					rt = append(rt, int(arr[j]))
				}
				return rt, true
			}
			if v, ok := i.(string); ok && v == "__DSL-ZERO__" {
				return []int{}, true
			}
			arr2, ok := i.([]interface{})
			if !ok {
				return nil, false
			}

			for j := 0; j < len(arr2); j++ {
				toVal, ok := convert.Int(arr2[j])
				if ok {
					rt = append(rt, toVal)
				}
			}
			if len(rt) == 0 && len(arr2) > 0 {
				return nil, false
			}

			//json-number
			arr3, ok := i.([]json.Number)
			if !ok {
				return nil, false
			}
			for j := 0; j < len(arr3); j++ {
				toVal, ok := convert.Int(arr2[j])
				if ok {
					rt = append(rt, toVal)
				}
			}
			if len(rt) == 0 && len(arr3) > 0 {
				return nil, false
			}
			return rt, true
		},

		consts.ARR_UINT32: func(i interface{}) (interface{}, bool) {
			rt := []int{}
			arr, ok := i.([]uint32)
			if ok {
				for j := 0; j < len(arr); j++ {
					rt = append(rt, int(arr[j]))
				}
				return rt, true
			}
			if v, ok := i.(string); ok && v == "__DSL-ZERO__" {
				return []int{}, true
			}
			arr2, ok := i.([]interface{})
			if !ok {
				return nil, false
			}

			for j := 0; j < len(arr2); j++ {
				toVal, ok := convert.Int(arr2[j])
				if ok {
					rt = append(rt, toVal)
				}
			}
			if len(rt) == 0 && len(arr2) > 0 {
				return nil, false
			}

			//json-number
			arr3, ok := i.([]json.Number)
			if !ok {
				return nil, false
			}
			for j := 0; j < len(arr3); j++ {
				toVal, ok := convert.Int(arr2[j])
				if ok {
					rt = append(rt, toVal)
				}
			}
			if len(rt) == 0 && len(arr3) > 0 {
				return nil, false
			}
			return rt, true
		},

		consts.ARR_UINT64: func(i interface{}) (interface{}, bool) {
			rt := []int{}
			arr, ok := i.([]uint64)
			if ok {
				for j := 0; j < len(arr); j++ {
					rt = append(rt, int(arr[j]))
				}
				return rt, true
			}
			if v, ok := i.(string); ok && v == "__DSL-ZERO__" {
				return []int{}, true
			}
			arr2, ok := i.([]interface{})
			if !ok {
				return nil, false
			}

			for j := 0; j < len(arr2); j++ {
				toVal, ok := convert.Int(arr2[j])
				if ok {
					rt = append(rt, toVal)
				}
			}
			if len(rt) == 0 && len(arr2) > 0 {
				return nil, false
			}

			//json-number
			arr3, ok := i.([]json.Number)
			if !ok {
				return nil, false
			}
			for j := 0; j < len(arr3); j++ {
				toVal, ok := convert.Int(arr2[j])
				if ok {
					rt = append(rt, toVal)
				}
			}
			if len(rt) == 0 && len(arr3) > 0 {
				return nil, false
			}
			return rt, true
		},

		consts.ARR_FLOAT32: func(i interface{}) (interface{}, bool) {
			rt := []float64{}
			arr, ok := i.([]float32)
			if ok {
				for j := 0; j < len(arr); j++ {
					rt = append(rt, float64(arr[j]))
				}
				return rt, true
			}
			if v, ok := i.(string); ok && v == "__DSL-ZERO__" {
				return []int{}, true
			}
			arr2, ok := i.([]interface{})
			if !ok {
				return nil, false
			}

			for j := 0; j < len(arr2); j++ {
				toVal, ok := convert.Float64(arr2[j])
				if ok {
					rt = append(rt, toVal)
				}
			}
			if len(rt) == 0 && len(arr2) > 0 {
				return nil, false
			}

			//json-number
			arr3, ok := i.([]json.Number)
			if !ok {
				return nil, false
			}
			for j := 0; j < len(arr3); j++ {
				toVal, ok := convert.Float64(arr2[j])
				if ok {
					rt = append(rt, toVal)
				}
			}
			if len(rt) == 0 && len(arr3) > 0 {
				return nil, false
			}
			return rt, true
		},

		consts.ARR_JSON_NUMBER: func(i interface{}) (interface{}, bool) {
			if v, ok := i.(string); ok && v == "__DSL-ZERO__" {
				return []float64{}, true
			}
			arr, ok := i.([]json.Number)
			if ok {
				if len(arr) == 0 {
					return []float64{}, true
				}

				if _, err := arr[0].Int64(); err == nil {
					rt := []int{}
					for j := 0; j < len(arr); j++ {
						val, err := arr[j].Int64()
						if err == nil {
							rt = append(rt, int(val))
						}
					}
					return rt, true
				} else {
					rt := []float64{}
					for j := 0; j < len(arr); j++ {
						val, err := arr[j].Float64()
						if err == nil {
							rt = append(rt, val)
						}
					}
					return rt, true
				}
			}

			arr2, ok := i.([]interface{})
			if !ok {
				return nil, false
			}
			rt := []float64{}
			for j := 0; j < len(arr2); j++ {
				toVal, ok := convert.Float64(arr2[j])
				if ok {
					rt = append(rt, toVal)
				}
			}
			if len(rt) == 0 && len(arr2) > 0 {
				return nil, false
			}
			return rt, true
		},

		consts.CONTINUE: func(i interface{}) (interface{}, bool) {
			is, _ := convert.Bool(i)
			if is {
				return "continue", true
			}
			return "", true
		},

		consts.BREAK: func(i interface{}) (interface{}, bool) {
			is, _ := convert.Bool(i)
			if is {
				return "break", true
			}
			return "", true
		},
	}

}

var once sync.Once

func init() {
	once.Do(func() {
		//注册
		initCvt()
	})
}
