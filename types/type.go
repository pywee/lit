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
