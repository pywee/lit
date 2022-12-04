package types

// 支持的类型
const (
	NULL      = "NULL"
	INTERFACE = "INTERFACE"
	STRING    = "STRING"
	INT       = "INT"
	FLOAT     = "FLOAT"
	BOOL      = "BOOL"
	FUNCTION  = "FUNC"
	ARRAY     = "ARRAY"
	OBJECT    = "OBJECT"
)

const (
	// 用于代码解析时区分代码块
	// IfExpressionNameIfTok 条件语句 if
	IfExpressionNameIfTok = "if"
	// IfExpressionNameElseIfTok 条件语句 elseif
	IfExpressionNameElseIfTok = "elseif"
	// IfExpressionNameElseTok 条件语句 "else"
	IfExpressionNameElseTok = "else"
)

const (
	// CodeTypeUnknow 0.未知, 非法定义
	CodeTypeUnknow = 0
	// CodeTypeIdentVAR 10.变量赋值
	CodeTypeIdentVAR = 10
	// CodeTypeIdentArrayVAR 11.为数组的指定下标赋值 可以是一个多维数组
	CodeTypeIdentArrayVAR = 11
	// CodeTypeIdentFN 20.函数声明
	CodeTypeIdentFN = 20
	// CodeTypeFunctionExec 21.函数调用
	CodeTypeFunctionExec = 21
	//CodeTypeIdentIF 30.if语句
	CodeTypeIdentIF = 30
	// CodeTypeIdentFOR 40.for
	CodeTypeIdentFOR = 40
	// CodeTypeIdentRETURN 50.return
	CodeTypeIdentRETURN = 50
	// CodeTypeVariablePlus 60.变量自增 n++
	CodeTypeVariablePlus = 60
	// CodeTypeVariableReduce 61.变量自减 n--
	CodeTypeVariableReduce = 61
	// CodeTypeContinue 70.continue
	CodeTypeContinue = 70
	// CodeTypeBreak 71.break
	CodeTypeBreak = 71
)

const (
	// TypeForExpressionIteration 1 for 循环迭代
	TypeForExpressionIteration = 1
	// TypeForExpressionRange 2 for range
	TypeForExpressionRange = 2
	// TypeForExpressionEndlessLoop 3 for死循环
	TypeForExpressionEndlessLoop = 3
)
