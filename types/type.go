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
	// CodeTypeVarIdent 1.变量赋值
	CodeTypeVarIdent = 1
	// CodeTypeFunctionIdent 2.函数声明
	CodeTypeFunctionIdent = 2
	// CodeTypeFunctionExec 3.函数调用
	CodeTypeFunctionExec = 3
	//CodeTypeIfIdent 4.if语句
	CodeTypeIfIdent = 4
	// CodeTypeForIdent 5.for
	CodeTypeForIdent = 5
	// CodeTypeReturnIdent 6.return
	CodeTypeReturnIdent = 6
)
