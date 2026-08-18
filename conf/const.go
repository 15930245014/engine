package conf

import "gitlab.shudieds.com/mec/lib/consts"

/*
*
token类型
*/
const (
	//函数
	Identifier = iota

	//数字
	Literal

	//运算符
	Operator

	// 通用的,
	COMMA

	//变量
	VARIABLE

	//输入
	INPUT

	//字段
	FIELD

	//全局
	GLOBAL
)

/*
*
计算类型
*/
const (
	//手填
	CalculateInput = "input"
	//函数获取
	CalculateFunc = "func"

	//固定值
	CalculateFixed = "fixed"

	//公式
	CalculateExpr = "expr"

	//map
	CalculateMap = "map"

	//array
	CalculateArray = "arr"
)

/*
*
返回值类型
*/
const (
	//---标准返回数据类型--
	ReturnInt     = consts.INT
	ReturnFloat   = consts.FLOAT
	ReturnStr     = consts.STR
	ReturnBool    = consts.BOOL
	ReturnAny     = consts.ANY
	ReturnMap     = consts.MAP
	ReturnArray   = consts.ARR
	ReturnDate    = consts.DATE
	ReturnDecimal = consts.DECIMAL

	//----全量支持数据类型---
	ReturnInt8    = consts.INT8
	ReturnInt16   = consts.INT16
	ReturnInt32   = consts.INT32
	ReturnInt64   = consts.INT64
	ReturnFloat32 = consts.FLOAT32
	ReturnUInt8   = consts.UINT8
	ReturnUInt16  = consts.UINT16
	ReturnUInt32  = consts.UINT32
	ReturnUInt64  = consts.UINT64
	ReturnUInt    = consts.UINT
	ReturnNil     = consts.NIL
	ReturnZero    = consts.ZERO
	ReturnSelf    = "self"

	//流程控制
	ReturnContinue = consts.CONTINUE
	ReturnBreak    = consts.BREAK
	ReturnArrStr   = consts.ARR_STR
)

// 优先级
var Precedence = map[string]int{
	"+":  30,
	"-":  30,
	"*":  40,
	"/":  40,
	"%":  40,
	"^":  60,
	"&&": 10,
	"&":  10,
	"|":  10,
	"||": 10,
	"!":  70,
	">":  20,
	">=": 20,
	"<":  20,
	"<=": 20,
	"==": 20,
	"!=": 20,
	"->": 80,
}

// 系统默认的常量字段
const SYS_NULL = "null"
const SYS_TRUE = "true"
const SYS_FALSE = "false"
const SYS_ZERO = "zero"
const SYS_SELF = "self"

// E 系统参数
const SYS_E = "E"

// I 系统参数
const SYS_I = "I"

// V 系统参数
const SYS_V = "V"

// 系统参数
const SYS_S = "S"

// R 公共参数
const SYS_R = "R"

// O 系统参数
const SYS_O = "O"

// 主数据
const SYS_M = "M"

// 对照
const SYS_L = "L"

// 类对象
const SYS_C = "C"

// OP 参数
const SYS_OP = "OP"

// OPP 参数
const SYS_OPP = "OPP"

// 获取外部传入
const SYS_ENV = "ENV"

// 主数据名称
const M_STORE = "store"
const M_MATERIALS = "materials"
const M_PLATFORM = "platform"

// dsl定义的系统变量
var SysSet = map[string]bool{
	SYS_I:   true,
	SYS_E:   true,
	SYS_V:   true,
	SYS_R:   true,
	SYS_C:   true,
	SYS_O:   true,
	SYS_OP:  true,
	SYS_OPP: true,
	SYS_M:   true,
	SYS_L:   true,
	SYS_S:   true,

	SYS_ENV:   true,
	SYS_NULL:  true,
	SYS_TRUE:  true,
	SYS_FALSE: true,
	SYS_ZERO:  true,
}

/*
*
* ast语法树定义的局部变量
 */

var SysLocalSet = map[string]bool{
	SYS_I: true,
	SYS_E: true,
	SYS_V: true,
}

/**
二级全局变量
*/

var SysSecondSet = map[string]bool{
	SYS_R:   true,
	SYS_S:   true,
	SYS_O:   true,
	SYS_OP:  true,
	SYS_OPP: true,
}

const CLASS_DWD = "dwd"
const CLASS_BILL = "bill"
const CLASS_DWS = "dws"
const CLASS_COMB = "comb"
