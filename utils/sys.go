package utils

import (
	"gitlab.shudieds.com/mec/lib/entry/engine"
	"gitlab.shudieds.com/zxh/engine/conf"
)

func GetOp(params map[string]interface{}) (interface{}, bool) {
	if params == nil {
		return nil, false
	}
	data, ok := params[conf.SYS_O]
	if !ok {
		return nil, false
	}
	O, ok := data.(*engine.TreeParams)
	if ok == false {
		return nil, false
	}

	return O.P, true
}
func GetOpp(params map[string]interface{}) (interface{}, bool) {
	OP, bl := GetOp(params)
	if !bl {
		return OP, false
	}
	return GetOp(map[string]interface{}{
		conf.SYS_O: OP,
	})
}
