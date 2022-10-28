package function

const (
	// baseFunctions
	// 支持的内置函数: 通用处理函数
	FUNCTION_PRINT   = "print"
	FUNCTION_VARDUMP = "varDump"

	// strFunctions
	// 支持的内置函数: 字符串处理函数
	FUNCTION_TRIM      = "trim"
	FUNCTION_LTRIM     = "trimLeft"
	FUNCTION_RTRIM     = "trimRight"
	FUNCTION_TRIMSPACE = "trimSpace"
	FUNCTION_LEN       = "len"
	FUNCTION_UTF8LEN   = "utf8Len"
	FUNCTION_MD5       = "md5"
	FUNCTION_REPLACE   = "replace"
	FUNCTION_CONTAINS  = "contains"
	FUNCTION_INDEX     = "index"
	FUNCTION_LASTINDEX = "lastIndex"
	FUNCTION_TOLOWER   = "toLower"
	FUNCTION_TOUPPER   = "toUpper"
	FUNCTION_TOTITLE   = "toTitle"
	FUNCTION_REPEAT    = "repeat"

	// 留到后面实现
	// FUNCTION_SPLIT      = "split"
	FUNCTION_SUBSTR = "substr"
	// FUNCTION_PARSEINT   = "parseInt"
	// FUNCTION_PARSEFLOAT = "parseFloat"

	// numberFunctions
	// 支持的内置函数: 数字处理函数
	FUNCTION_INT       = "int"
	FUNCTION_FLOOR     = "floor"
	FUNCTION_STRING    = "string"
	FUNCTION_ISNUMERIC = "isNumeric"
	FUNCTION_ISBOOL    = "isBool"
	FUNCTION_ISINT     = "isInt"
	FUNCTION_ISFLOAT   = "isFloat"
)
