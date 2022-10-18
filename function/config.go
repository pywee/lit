package function

const (

	// baseFunctions
	// 支持的内置函数: 通用处理函数
	FUNCTION_PRINT   = "Print"
	FUNCTION_VARDUMP = "VarDump"

	// strFunctions
	// 支持的内置函数: 字符串处理函数
	FUNCTION_TRIM      = "Trim"
	FUNCTION_LTRIM     = "TrimLeft"
	FUNCTION_RTRIM     = "TrimRight"
	FUNCTION_TRIMSPACE = "TrimSpace"
	FUNCTION_LEN       = "Len"
	FUNCTION_UTF8LEN   = "UTF8Len"
	FUNCTION_MD5       = "MD5"
	FUNCTION_REPLACE   = "Replace"
	FUNCTION_CONTAINS  = "Contains"
	FUNCTION_INDEX     = "Index"
	FUNCTION_LASTINDEX = "LastIndex"
	FUNCTION_TOLOWER   = "ToLower"
	FUNCTION_TOUPPER   = "ToUpper"
	FUNCTION_TOTITLE   = "ToTitle"
	FUNCTION_REPEAT    = "Repeat"

	// 留到后面实现
	// FUNCTION_SPLIT      = "split"
	FUNCTION_SUBSTR = "Substr"
	// FUNCTION_PARSEINT   = "parseInt"
	// FUNCTION_PARSEFLOAT = "parseFloat"

	// numberFunctions
	// 支持的内置函数: 数字处理函数
	FUNCTION_INT       = "Int"
	FUNCTION_FLOOR     = "Floor"
	FUNCTION_STRING    = "String"
	FUNCTION_ISNUMERIC = "IsNumeric"
	FUNCTION_ISBOOL    = "IsBool"
	FUNCTION_ISINT     = "IsInt"
	FUNCTION_ISFLOAT   = "IsFloat"
)
