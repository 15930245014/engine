package dsl

import (
	"encoding/json"
	"errors"
	"github.com/iancoleman/orderedmap"
	"gitlab.shudieds.com/mec/lib/utils/uuid"
	"gitlab.shudieds.com/zxh/engine/utils"
	"math"
	"math/rand"
	"sort"
	"strings"
	"time"

	"gitlab.shudieds.com/mec/lib/utils/convert"
	"gitlab.shudieds.com/mec/lib/utils/jsoniter"
	"gitlab.shudieds.com/mec/lib/utils/maphelper"
	"gitlab.shudieds.com/zxh/engine/types"
)

const (
	RadianMode = iota
	AngleMode
)

type defS struct {
	argc int
	fun  func(expr ...types.ExprAST) (interface{}, error)
}

// enum "RadianMode", "AngleMode"
var TrigonometricMode = RadianMode
var tokenCheckDef = map[string]bool{
	//数字相关
	"sin":   true,
	"cos":   true,
	"tan":   true,
	"cot":   true,
	"sec":   true,
	"csc":   true,
	"abs":   true,
	"ceil":  true,
	"floor": true,
	"round": true,
	"sqrt":  true,
	"cbrt":  true,
	"max":   true,
	"min":   true,

	//字符串相关
	"strLen":        true,
	"strConcat":     true,
	"strIn":         true,
	"strNotIn":      true,
	"strToUpper":    true,
	"strToLower":    true,
	"strHasPrefix":  true,
	"strHasSuffix":  true,
	"strContains":   true,
	"strCount":      true,
	"strIndex":      true,
	"strLastIndex":  true,
	"strTrimSpace":  true,
	"strTrim":       true,
	"strTrimLeft":   true,
	"strTrimRight":  true,
	"strTrimPrefix": true,
	"strTrimSuffix": true,
	"strReplace":    true,
	"strSplit":      true,
	"strJoin":       true,
	"strRepeat":     true,
	"strConcatFill": true,
	"strSub":        true,
	"strLastCut":    true,
	"strStr":        true,

	//加密相关
	"base64Encode": true,
	"base64Decode": true,
	"md5":          true,
	"hmac":         true,
	"hmacSha256":   true,
	"sha1":         true,
	"sha256":       true,
	"encrypt":      true,
	"aesECBEncode": true,
	"aesECBDecode": true,
	"aesCBCEncode": true,
	"aesCBCDecode": true,
	"aesCRTEncode": true,
	"aesCRTDecode": true,
	"aesCFBEncode": true,
	"aesCFBDecode": true,
	"aesOFBEncode": true,
	"aesOFBDecode": true,

	//json
	"jsonEncode": true,
	"jsonDecode": true,

	//流程控制
	"if":      true,
	"switch":  true,
	"choice":  true,
	"case":    true,
	"default": true,
	"foreach": true,
	"each":    true,

	//数组切片函数
	"arrLen":      true,
	"arrIndex":    true,
	"arrFirst":    true,
	"arrLast":     true,
	"arrAppend":   true,
	"arrUnshift":  true,
	"arrSlice":    true,
	"arrMerge":    true,
	"arrSearch":   true,
	"arrUnique":   true,
	"arrReverse":  true,
	"arrSortAsc":  true,
	"arrSortDesc": true,

	//map函数
	"mapLen":      true,
	"mapExist":    true,
	"mapGet":      true,
	"mapSet":      true,
	"mapKeys":     true,
	"mapValues":   true,
	"mapMerge":    true, //如果不存在 在添加
	"mapSortAsc":  true,
	"mapSortDesc": true,
	"mapRef":      true,
	"mapIsEmpty":  true,
	"mapToStr":    true,

	//时间相关
	"dateTimestamp":       true,
	"dateFormat":          true,
	"dateGetYear":         true,
	"dateGetMonth":        true,
	"dateGetDay":          true,
	"dateGetHour":         true,
	"dateGetMinute":       true,
	"dateGetSecond":       true,
	"dateGetLastMonthDay": true,
	"getTimestamp":        true,
	"getCurrentTime":      true,
	"getYYmm":             true,
	"getUnixMilliTime":    true,

	//类型判断与类型转化
	"isInt":     true,
	"isFloat":   true,
	"isStr":     true,
	"isDate":    true,
	"isBool":    true,
	"isArray":   true,
	"isObj":     true,
	"isDecimal": true,
	"isZero":    true,

	//类型转化
	"toInt":     true,
	"toFloat":   true,
	"toStr":     true,
	"toDate":    true,
	"toBool":    true,
	"toDecimal": true,

	//比较函数
	"compare": true,
	"equal":   true,

	//自定函数
	"getJJToken":      true,
	"getWDTAuth":      true,
	"getJSTsign":      true,
	"getJDYSessionId": true,
	"getEccangSign":   true,
	"getWDTQMSign":    true,
	"getWDTSign":      true,
	"getSellfoxSign":  true,
	"getSellfoxToken": true,
	"getKDNSign":      true,
	"getPDDSign":      true,

	//输入
	"env": true,

	//逻辑运算符
	"and": true,
	"or":  true,
	"not": true,

	//常量
	"null":  true,
	"true":  true,
	"false": true,
	"zero":  true,
	"self":  true,

	//其他语言代码
	"jsExec": true,

	//decimal
	"decimalAvg":    true, //求多个平均数
	"decimalMax":    true, //求多个最大值
	"decimalMin":    true, //求多个最小值
	"decimalAdd":    true, //求多个的和
	"decimalSub":    true, //减法
	"decimalMul":    true, //相乘
	"decimalDiv":    true, //除数相除
	"decimalAbs":    true, //求平均数
	"decimalCeil":   true, //向上取整
	"decimalFloor":  true, //向上取整
	"decimalCmp":    true, //比较大小-1,0,1
	"decimalEqual":  true, //是否相等
	"decimalIsZero": true, //是否为0
	"decimalLt":     true, //是否小于
	"decimalRt":     true, //是否大于等于
	"decimalRte":    true, //是否大于等于
	"decimalLte":    true, //是否小于等于
	"decimalMod":    true, //取模
	"decimalNeg":    true, //取反
	"decimalPow":    true, //取方
	"decimalRound":  true, //四舍五入，位数

	//业务主数据
	"MGet":   true,
	"LGet":   true,
	"LGetV2": true,

	//位运算
	"bit":  true,
	"expr": true,

	//try
	"try": true,

	//随机数
	"randInt": true,

	//父函数
	"P": true,

	/*
	  类函数 ---------------------------------
	*/
	"getEsDwdInformation":      true,
	"getEsDwdChildInformation": true,
	"getEsBillM2":              true,
	"getEsDwsStatement":        true,
	"getEsComb":                true,

	//uuid
	"uuidV3": true,
}

func (a *AST) initFunc() {
	a.defFunc = map[string]defS{
		//数字相关
		"sin":  {1, a.defSin},
		"cos":  {1, a.defCos},
		"tan":  {1, a.defTan},
		"cot":  {1, a.defCot},
		"sec":  {1, a.defSec},
		"csc":  {1, a.defCsc},
		"cbrt": {1, a.defCbrt},

		"abs":   {1, a.defAbs},
		"ceil":  {1, a.defCeil},
		"floor": {1, a.defFloor},
		"round": {2, a.defRound},
		"sqrt":  {1, a.defSqrt},

		"max": {-1, a.defMax},
		"min": {-1, a.defMin},

		//字符串相关
		"strLen":        {1, a.strLen},
		"strConcat":     {-1, a.strConcat},
		"strIn":         {-1, a.strIn},
		"strNotIn":      {-1, a.strNotIn},
		"strToUpper":    {1, a.strToUpper},
		"strToLower":    {1, a.strToLower},
		"strHasPrefix":  {2, a.strHasPrefix},
		"strHasSuffix":  {2, a.strHasSuffix},
		"strContains":   {2, a.strContains},
		"strCount":      {2, a.strCount},
		"strIndex":      {2, a.strIndex},
		"strLastIndex":  {2, a.strLastIndex},
		"strTrimSpace":  {1, a.strTrimSpace},
		"strTrim":       {2, a.strTrim},
		"strTrimLeft":   {2, a.strTrimLeft},
		"strTrimRight":  {2, a.strTrimRight},
		"strTrimPrefix": {2, a.strTrimPrefix},
		"strTrimSuffix": {2, a.strTrimSuffix},
		"strReplace":    {3, a.strReplace},
		"strSplit":      {2, a.strSplit},
		"strJoin":       {2, a.strJoin},
		"strRepeat":     {2, a.strRepeat},
		"strConcatFill": {-1, a.strConcatFill},
		"strStr":        {1, a.strStr},
		"strSub":        {3, a.strSub},
		"strLastCut":    {2, a.strLastCut},

		//加密相关
		"base64Encode": {1, a.strBase64Encode},
		"base64Decode": {1, a.strBase64Decode},
		"md5":          {1, a.strMd5},
		"hmac":         {2, a.strHmac},
		"hmacSha256":   {2, a.strHmacSha256},
		"sha1":         {1, a.strSha1},
		"sha256":       {1, a.strSha256},
		"encrypt":      {1, a.strEncrypt},
		"aesECBEncode": {2, a.strAesECBEncode},
		"aesECBDecode": {2, a.strAesECBDecode},
		"aesCBCEncode": {2, a.strAesCBCEncode},
		"aesCBCDecode": {2, a.strAesCBCDecode},
		"aesCRTEncode": {2, a.strAesCRTEncode},
		"aesCRTDecode": {2, a.strAesCRTDecode},
		"aesCFBEncode": {2, a.strAesCFBEncode},
		"aesCFBDecode": {2, a.strAesCFBDecode},
		"aesOFBEncode": {2, a.strAesOFBEncode},
		"aesOFBDecode": {2, a.strAesOFBDecode},

		//json
		"jsonEncode": {-1, a.defJsonEncode},
		"jsonDecode": {1, a.defJsonDecode},

		//流程控制
		"if":      {3, a.defIf},      //if-else
		"switch":  {-1, a.defSwitch}, //switch(值,case(条件1,do),case(条件2，do),case(条件3，do),default(默认执行))
		"choice":  {-1, a.defChoice}, //chose(case(条件1,do),case(条件2，do),case(条件3，do),default(默认执行))
		"case":    {2, a.defCase},    //switch(case1,do1,case2,do2)
		"default": {1, a.defDefault}, //switch(case1,do1,case2,do2)

		"foreach": {-1, a.defForeach},
		"each":    {4, a.defEach},

		//逻辑运算符
		"and": {-1, a.defAnd},
		"or":  {-1, a.defOr},
		"not": {1, a.defNot},

		//数组切片函数
		"arrLen":      {1, a.defArrLen},
		"arrIndex":    {2, a.defArrIndex},
		"arrFirst":    {1, a.arrFirst},
		"arrLast":     {1, a.arrLast},
		"arrAppend":   {2, a.arrAppend},
		"arrUnshift":  {2, a.arrUnshift},
		"arrSlice":    {3, a.arrSlice},
		"arrMerge":    {2, a.arrMerge},
		"arrSearch":   {2, a.arrSearch},
		"arrUnique":   {1, a.arrUnique},
		"arrReverse":  {1, a.arrReverse},
		"arrSortAsc":  {1, a.arrSortAsc},
		"arrSortDesc": {1, a.arrSortDesc},

		//map函数
		"mapLen":      {1, a.mapLen},
		"mapIsEmpty":  {1, a.mapIsEmpty},
		"mapExist":    {2, a.mapKeyExist},
		"mapGet":      {-1, a.mapGet},
		"mapSet":      {3, a.mapSet},
		"mapKeys":     {1, a.mapKeys},
		"mapValues":   {1, a.mapValues},
		"mapMerge":    {2, a.mapMerge}, //如果不存在 在添加
		"mapSortAsc":  {1, a.mapKeySortAsc},
		"mapSortDesc": {1, a.mapKeySortDesc},
		"mapRef":      {1, a.mapRef},
		"mapToStr":    {1, a.mapToStr},

		//时间相关
		"dateTimestamp":       {1, a.dateTimestamp},
		"dateFormat":          {2, a.dateTimeFormat},
		"dateGetYear":         {1, a.dateGetYear},
		"dateGetMonth":        {1, a.dateGetMonth},
		"dateGetDay":          {1, a.dateGetDay},
		"dateGetHour":         {1, a.dateGetHour},
		"dateGetMinute":       {1, a.dateGetMinute},
		"dateGetSecond":       {1, a.dateGetSecond},
		"dateGetLastMonthDay": {1, a.dateGetLastMonthDay},
		"getTimestamp":        {0, a.getTimestamp},
		"getCurrentTime":      {0, a.getCurrentTime},
		"getYYmm":             {3, a.getYYmm},
		"getUnixMilliTime":    {0, a.getUnixMilliTime},

		//类型判断函数
		"isInt":     {1, a.isInt},
		"isFloat":   {1, a.isFloat},
		"isStr":     {1, a.isStr},
		"isDate":    {1, a.isDate},
		"isBool":    {1, a.isBool},
		"isArray":   {1, a.isArray},
		"isObj":     {1, a.isObj},
		"isDecimal": {1, a.isDecimal},
		"isZero":    {1, a.isZero},

		//转换函数
		"toInt":     {-1, a.toInt},
		"toFloat":   {-1, a.toFloat},
		"toStr":     {-1, a.toStr},
		"toDate":    {-1, a.toDate},
		"toBool":    {-1, a.toBool},
		"toDecimal": {-1, a.toDecimal},

		//比较函数
		"compare": {2, a.defCompare},
		"equal":   {2, a.defEqual},

		//自定义函数
		"getToken":        {3, a.getJJToken},
		"getWDTAuth":      {5, a.getWDTAuth},
		"getJSTsign":      {5, a.getJSTsign},
		"getJDYSessionId": {0, a.getJDYSessionId},
		"getEccangSign":   {8, a.getEccangSign},
		"getWDTQMSign":    {1, a.getWDTQMSign},
		"getWDTSign":      {3, a.getWDTSign},
		"getSellfoxSign":  {7, a.getSellfoxSign},
		"getSellfoxToken": {2, a.getSellfoxToken},
		"getPDDSign":      {2, a.getPDDSign},

		//执行其他语言代码
		"jsExec": {-1, a.jsExec},

		//decimal
		"decimalAvg":    {-1, a.decimalAvg},   //求多个平均数
		"decimalMax":    {-1, a.decimalMax},   //求多个最大值
		"decimalMin":    {-1, a.decimalMin},   //求多个最小值
		"decimalAdd":    {-1, a.decimalAdd},   //求多个的和
		"decimalSub":    {-1, a.decimalSub},   //减法
		"decimalMul":    {-1, a.decimalMul},   //相乘
		"decimalDiv":    {-1, a.decimalDiv},   //除数相除
		"decimalAbs":    {1, a.decimalAbs},    //求平均数
		"decimalCeil":   {1, a.decimalCeil},   //向上取整
		"decimalFloor":  {1, a.decimalFloor},  //向下	取整
		"decimalCmp":    {2, a.decimalCmp},    //比较大小-1,0,1
		"decimalEqual":  {2, a.decimalEqual},  //是否相等
		"decimalIsZero": {1, a.decimalIsZero}, //是否为0
		"decimalRt":     {2, a.decimalRt},     //是否大于
		"decimalRte":    {2, a.decimalRte},    //是否大于等于
		"decimalLt":     {2, a.decimalLt},     //是否小于
		"decimalLte":    {2, a.decimalLte},    //是否小于等于
		"decimalMod":    {2, a.decimalMod},    //取模
		"decimalNeg":    {1, a.decimalNeg},    //取反
		"decimalPow":    {2, a.decimalPow},    //取方
		"decimalRound":  {2, a.decimalRound},  //四舍五入，位数

		//主数据
		"MGet":   {1, a.MGet},
		"LGet":   {3, a.LGet},
		"LGetV2": {4, a.LGetV2},
		"P":      {1, a.P},

		//位运算
		"bit": {-1, a.bit},

		//expr
		"expr": {1, a.defExpr},

		//try
		"try": {2, a.defTry},

		//随机数
		"randInt": {2, a.defRandInt},

		/*
		  类函数 ---------------------------------
		*/
		//dwd
		"getEsDwdInformation":      {5, a.GetEsDwdInformation},
		"getEsDwdChildInformation": {5, a.GetEsDwdChildInformation},
		//bill
		"getEsBillM2": {4, a.GetEsBillM2},
		//dws
		"getEsDwsStatement": {1, a.GetEsDwsStatement},
		"getEsComb":         {5, a.GetEsComb},

		//uuid
		"uuidV3": {1, a.defUuidV3},
	}
}

func (a *AST) defUuidV3(exprs ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(exprs[0])
	if err != nil {
		return nil, errors.New("uuidV3函数执行错误，参数1解析错误")
	}
	s, bl := convert.String(v1)
	if bl == false {
		return nil, errors.New("uuidV3函数执行错误，参数1应为string")
	}

	return uuid.GenUuidV3(s), nil
}

/*
*
env
*/

func (a *AST) defExpr(exprs ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(exprs[0])
	if err != nil {
		return nil, err
	}
	ex, ok := convert.String(v)
	if ok == false || ex == "" {
		return nil, errors.New("expr 表达式错误")
	}
	if a.exprCal == nil {
		a.exprCal = NewExprCalculate(a.c)
		a.exprCal.SetFieldCalculate(a.fieldCalculate)
		a.exprCal.SetVariableCalculate(a.variableCalculate)
	}
	a.exprCal.Set(ex, ex, a.globalParams, ex)
	return a.exprCal.Calculate()
}

/*
*
env
*/

func (a *AST) defTry(exprs ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(exprs[0])
	if err != nil || utils.IsNil(v) {
		v2, err := a.ExprASTResult(exprs[1])
		if err != nil {
			return nil, errors.New("try函数执行错误，前后参数全部错误!")
		}
		return v2, nil
	} else {
		return v, nil
	}

}

func (a *AST) defRandInt(exprs ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(exprs[0])
	if err != nil {
		return nil, errors.New("randInt函数执行错误，参数1解析错误")
	}
	s, bl := convert.Int(v1)
	if bl == false {
		return nil, errors.New("randInt函数执行错误，参数1应为int")
	}

	v2, err := a.ExprASTResult(exprs[1])
	if err != nil {
		return nil, errors.New("randInt函数执行错误，参数2解析错误")
	}

	e, bl := convert.Int(v2)
	if bl == false {
		return nil, errors.New("randInt函数执行错误，参数2应为int")
	}

	if s >= e {
		return nil, errors.New("randInt函数执行错误，参数1应小于参数2")
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(e-s+1) + s, nil
}

/*
*
env
*/
func (a *AST) defJs(exprs ...types.ExprAST) (interface{}, error) {
	for _, expr := range exprs {
		v, err := a.ExprASTResult(expr)
		if err != nil {
			return nil, err
		}
		bl, ok := convert.Bool(v)
		if ok == false {
			return nil, errors.New("转化bool类型失败")
		}
		if bl == true {
			return true, nil
		}
	}
	return false, nil
}

/*
*not
 */
func (a *AST) defNot(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	bl, ok := convert.Bool(v)
	if ok == false {
		return nil, errors.New("not函数参数返回应该为bool类型!")
	}

	return !bl, nil
}

/*
*and
 */
func (a *AST) defAnd(expr ...types.ExprAST) (interface{}, error) {
	if len(expr) == 0 {
		return nil, errors.New("and函数参数不能为空!")
	}
	for i := 0; i < len(expr); i++ {
		v, err := a.ExprASTResult(expr[i])
		if err != nil {
			return nil, err
		}
		bl, ok := convert.Bool(v)
		if ok == false {
			return nil, errors.New("and 函数参数返回应该为bool类型!")
		}
		if bl == false {
			return false, nil
		}
	}
	return true, nil
}

/*
*or
 */
func (a *AST) defOr(expr ...types.ExprAST) (interface{}, error) {
	if len(expr) == 0 {
		return nil, errors.New("or 函数参数不能为空!")
	}
	for i := 0; i < len(expr); i++ {
		v, err := a.ExprASTResult(expr[i])
		if err != nil {
			return nil, err
		}
		bl, ok := convert.Bool(v)
		if ok == false {
			return nil, errors.New("or 函数参数返回应该为bool类型!")
		}
		if bl == true {
			return true, nil
		}
	}
	return false, nil
}

/*
*
env
*/
func (a *AST) defEnv(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	str, bl := convert.String(v)
	if bl == false {
		return nil, errors.New("env解析参数错误!")
	}

	val, bl := maphelper.JDGet(a.globalParams, strings.Split(str, "."))
	if bl == false {
		return nil, errors.New("env获取值失败")
	}
	return val, nil
}

func (a *AST) defJsonEncode(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	if len(expr) == 1 {
		//fmt.Println("kkk：", v)
		data, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		//fmt.Println(string(data))
		return convert.BytesToString(data), nil
	}

	toMp, ok := convert.ToMap(v)
	if !ok {
		data, err := jsoniter.Marshal(v)
		if err != nil {
			return "", err
		}
		return convert.BytesToString(data), nil
	}

	v1, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}

	isOrder, bl := convert.Bool(v1)
	if !bl {
		return nil, errors.New("jsonDecode 解析参数错误")
	}

	if !isOrder {
		data, err := jsoniter.Marshal(v)
		if err != nil {
			return "", err
		}
		return convert.BytesToString(data), nil
	}
	var keys []string
	for k, _ := range toMp {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	m := orderedmap.New()
	for i := 0; i < len(keys); i++ {
		m.Set(keys[i], toMp[keys[i]])
	}
	data, err := jsoniter.Marshal(m)
	if err != nil {
		return "", err
	}
	return convert.BytesToString(data), nil
}

/*
*
jsonDecode
*/
func (a *AST) defJsonDecode(expr ...types.ExprAST) (interface{}, error) {
	jsonStr, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	str, bl := convert.String(jsonStr)
	if !bl {
		return nil, errors.New("jsonDecode 解析参数错误")
	}

	//创建一个接口来保存解析后的数据
	var result map[string]interface{}
	err = json.Unmarshal(convert.StringToBytes(str), &result)
	if err != nil {
		return nil, errors.New("jsonDecode 解析参数错误:" + err.Error())
	}
	return result, nil

}

/****数字相关*****/
// sin(pi/2) = 1
func (a *AST) defSin(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.expr2Radian(expr[0])
	if err != nil {
		return nil, err
	}
	return math.Sin(v), nil
}

// cos(0) = 1
func (a *AST) defCos(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.expr2Radian(expr[0])
	if err != nil {
		return nil, err
	}
	return math.Cos(v), nil
}

// tan(pi/4) = 1
func (a *AST) defTan(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.expr2Radian(expr[0])
	if err != nil {
		return nil, err
	}

	return math.Tan(v), nil
}
func (a *AST) expr2Radian(expr types.ExprAST) (float64, error) {
	v, err := a.ExprASTResult(expr)
	if err != nil {
		return 0, err
	}

	r, _ := convert.Float64(v)
	if TrigonometricMode == AngleMode {
		r = r / 180 * math.Pi
	}
	return r, nil
}

// cot(pi/4) = 1
func (a *AST) defCot(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.defTan(expr...)
	if err != nil {
		return nil, err
	}
	f, _ := convert.Float64(v)
	if f == 0 {
		return nil, errors.New("defCot divisor is 0")
	}
	return 1 / f, nil
}

// sec(0) = 1
func (a *AST) defSec(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.defCos(expr...)
	if err != nil {
		return nil, err
	}
	f, _ := convert.Float64(v)
	if f == 0 {
		return nil, errors.New("defSec divisor is 0")
	}
	return 1 / f, nil
}

// csc(pi/2) = 1
func (a *AST) defCsc(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.defSin(expr...)
	if err != nil {
		return nil, err
	}
	f, _ := convert.Float64(v)
	if f == 0 {
		return nil, errors.New("defCsc divisor is 0")
	}
	return 1 / f, nil
}

// abs(-2) = 2
func (a *AST) defAbs(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	return convert.Abs(v)
}

// ceil(4.2) = ceil(4.8) = 5
func (a *AST) defCeil(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	return convert.Ceil(v)
}

// floor(4.2) = floor(4.8) = 4
func (a *AST) defFloor(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	return convert.Floor(v)
}

// round(4.2) = 4
// round(4.6) = 5
func (a *AST) defRound(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}
	dig, ok := convert.Int(v2)
	if !ok {
		return nil, errors.New("round 第二个参数应为整型")
	}
	return convert.Round(v, dig)
}

// sqrt(4) = 2
// sqrt(4) = abs(sqrt(4))
// returns only the absolute value of the result
func (a *AST) defSqrt(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	return convert.Sqrt(v)
}

// cbrt(27) = 3
func (a *AST) defCbrt(expr ...types.ExprAST) (interface{}, error) {
	v, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	f, _ := convert.Float64(v)
	return math.Cbrt(f), nil
}

// max(2) = 2
// max(2, 3) = 3
// max(2, 3, 1) = 3
func (a *AST) defMax(expr ...types.ExprAST) (interface{}, error) {
	if len(expr) == 0 {
		return nil, errors.New("calling function `max` must have at least one parameter.")
	}
	if len(expr) == 1 {
		v, err := a.ExprASTResult(expr[0])
		if err != nil {
			return nil, err
		}
		return v, nil
	}

	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	maxV := v1
	for i := 1; i < len(expr); i++ {
		v2, err := a.ExprASTResult(expr[i])
		if err != nil {
			return nil, err
		}
		res, err := convert.Compare(maxV, v2)
		if err != nil {
			return nil, errors.New("max参数错误:" + err.Error())
		}
		if res == -1 {
			maxV = v2
		}
	}
	return maxV, nil
}

// min(2) = 2
// min(2, 3) = 2
// min(2, 3, 1) = 1
func (a *AST) defMin(expr ...types.ExprAST) (interface{}, error) {
	if len(expr) == 0 {
		return nil, errors.New("calling function `min` must have at least one parameter.")
	}
	if len(expr) == 1 {
		v, err := a.ExprASTResult(expr[0])
		if err != nil {
			return nil, err
		}
		return v, nil
	}

	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	minV := v1
	for i := 1; i < len(expr); i++ {
		v2, err := a.ExprASTResult(expr[i])
		if err != nil {
			return nil, err
		}
		res, err := convert.Compare(minV, v2)
		if err != nil {
			return nil, errors.New("max参数错误:" + err.Error())
		}
		if res == 1 {
			minV = v2
		}
	}
	return minV, nil
}
