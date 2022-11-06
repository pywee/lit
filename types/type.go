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
	// CodeTypeIdentVAR 1.变量赋值
	CodeTypeIdentVAR = 1
	// CodeTypeIdentFN 2.函数声明
	CodeTypeIdentFN = 2
	// CodeTypeFunctionExec 3.函数调用
	CodeTypeFunctionExec = 3
	//CodeTypeIdentIF 4.if语句
	CodeTypeIdentIF = 4
	// CodeTypeIdentFOR 5.for
	CodeTypeIdentFOR = 5
	// CodeTypeIdentRETURN 6.return
	CodeTypeIdentRETURN = 6
)
