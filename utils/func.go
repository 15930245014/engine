package utils

import (
	"reflect"
	"strconv"
	"time"
)

// 判断接口是否为空
func IsNil(input interface{}) bool {
	if input == nil {
		return true
	}
	if reflect.TypeOf(input).Kind() == reflect.Ptr && reflect.ValueOf(input).IsNil() {
		return true
	}
	return false
}

func GetMonthStartAndEnd(myYear string, myMonth string) (string, string) {
	// 数字月份必须前置补零
	if len(myMonth) == 1 {
		myMonth = "0" + myMonth
	}
	yInt, _ := strconv.Atoi(myYear)
	theTime, _ := time.ParseInLocation("2006-01-02 15:04:05", myYear+"-"+myMonth+"-01 00:00:00", time.Local)
	newMonth := theTime.Month()

	t1 := time.Date(yInt, newMonth, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02")
	t2 := time.Date(yInt, newMonth+1, 0, 0, 0, 0, 0, time.Local).Format("2006-01-02")

	return t1 + " 00:00:00", t2 + " 23:59:59"
}
