package function

const (
	// baseFunctions
	// 支持的内置函数: 通用处理函数
	FUNCTION_PRINT     = "println"
	FUNCTION_VARDUMP   = "varDump"
	FUNCTION_ISBOOL    = "isBool"
	FUNCTION_ISINT     = "isInt"
	FUNCTION_ISFLOAT   = "isFloat"
	FUNCTION_ISNUMERIC = "isNumeric"
	FUNCTION_MD5       = "md5"

	// strFunctions
	// 支持的内置函数: 字符串处理函数

	// FUNCTION_TRIM
	// trim 返回去掉两边指定字符后的字符串 默认去掉空格
	// trim(str, [charlist=' ' string...]) string
	FUNCTION_TRIM = "trim"

	// FUNCTION_LTRIM
	// trimLeft 返回去掉左边指定字符后的字符串 默认去掉空格
	// trimLeft(str, [charlist='', string...]) string
	FUNCTION_LTRIM = "trimLeft"

	// FUNCTION_RTRIM
	// trimRight 返回去掉右边指定字符后的字符串 默认去掉空格
	// trimRight(str, [charlist='', string...]) string
	FUNCTION_RTRIM = "trimRight"

	// FUNCTION_TRIMSPACE
	// trimSpace 返回将 str 两端所有空白都去掉之后的字符串 包含空格 \n \r\n \t
	// trimSpace(str string) string
	FUNCTION_TRIMSPACE = "trimSpace"

	// FUNCTION_LEN
	// len 返回字符串长度或数组长度 (支持半角和全角长度判断)
	// len(str, [rune bool]) int
	FUNCTION_LEN = "len"

	FUNCTION_REPLACE    = "replace"
	FUNCTION_CONTAINS   = "contains"
	FUNCTION_INDEX      = "index"
	FUNCTION_LASTINDEX  = "lastIndex"
	FUNCTION_TOLOWER    = "toLower"
	FUNCTION_TOUPPER    = "toUpper"
	FUNCTION_TOTITLE    = "toTitle"
	FUNCTION_REPEAT     = "repeat"
	FUNCTION_SPLIT      = "split"
	FUNCTION_SUBSTR     = "substr"
	FUNCTION_PARSEINT   = "parseInt"
	FUNCTION_PARSEFLOAT = "parseFloat"

	// numberFunctions
	// 支持的内置函数: 数字处理函数
	FUNCTION_INT    = "int"
	FUNCTION_FLOOR  = "floor"
	FUNCTION_STRING = "string"
)
