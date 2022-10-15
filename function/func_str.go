package function

import (
	"strconv"
	"strings"

	"github.com/pywee/goExpr/global"
)

// strFunctions
// 支持的内置函数: 字符串处理函数
const (
	FUNCTION_REPLACE   = "replace"
	FUNCTION_TRIM      = "trim"
	FUNCTION_LTRIM     = "ltrim"
	FUNCTION_RTRIM     = "rtrim"
	FUNCTION_TRIMSPACE = "trimSpace"
	FUNCTION_PARSEINT  = "parseInt"
	FUNCTION_SPLIT     = "split"
	FUNCTION_MD5       = "md5"
)

// strFunctions
// 支持的内置函数: 字符串处理函数
var strFunctions = []*functionInfo{
	{
		FunctionName: FUNCTION_REPLACE,
		MustAmount:   1,
		MaxAmount:    4,
		Args: []*functionArgAttr{
			{Type: TYPE_INTERFACE, Must: true},
			{Type: TYPE_STRING, Must: true},
			{Type: TYPE_STRING, Must: true},
			{Type: TYPE_INT, Must: true},
		},
		FN: func(args ...*global.Structure) (*global.Structure, error) {
			a0 := args[0].Lit
			a1 := args[1].Lit
			a2 := args[2].Lit
			a3, _ := strconv.Atoi(args[3].Lit)
			rx := strings.Replace(a0, a1, a2, a3)
			return &global.Structure{Tok: TYPE_STRING, Lit: rx}, nil
		},
	},
	{
		// 必选参数1个
		// 可选参数1个 当不输入可选参数时 默认为空格 " "
		FunctionName: FUNCTION_TRIM,
		MustAmount:   1,
		MaxAmount:    2,
		Args: []*functionArgAttr{
			{Type: TYPE_INTERFACE, Must: true},
			{Type: TYPE_STRING, Must: false},
		},
		FN: func(args ...*global.Structure) (*global.Structure, error) {
			if len(args) > 1 {
				return &global.Structure{Tok: TYPE_STRING, Lit: strings.Trim(args[0].Lit, args[1].Lit)}, nil
			}
			return &global.Structure{Tok: TYPE_STRING, Lit: strings.Trim(args[0].Lit, " ")}, nil
		},
	},
	{
		// 必选参数1个
		// 可选参数1个 当不输入可选参数时 默认为空格 " "
		FunctionName: FUNCTION_LTRIM,
		MustAmount:   1,
		MaxAmount:    2,
		Args: []*functionArgAttr{
			{Type: TYPE_INTERFACE, Must: true},
			{Type: TYPE_STRING, Must: false},
		},
		FN: func(args ...*global.Structure) (*global.Structure, error) {
			if len(args) > 1 {
				return &global.Structure{Tok: TYPE_STRING, Lit: strings.TrimLeft(args[0].Lit, args[1].Lit)}, nil
			}
			return &global.Structure{Tok: TYPE_STRING, Lit: strings.TrimLeft(args[0].Lit, " ")}, nil
		},
	},
	{
		// 必选参数1个
		// 可选参数1个 当不输入可选参数时 默认为空格 " "
		FunctionName: FUNCTION_RTRIM,
		MustAmount:   1,
		MaxAmount:    2,
		Args: []*functionArgAttr{
			{Type: TYPE_INTERFACE, Must: true},
			{Type: TYPE_STRING, Must: false},
		},
		FN: func(args ...*global.Structure) (*global.Structure, error) {
			if len(args) > 1 {
				return &global.Structure{Tok: TYPE_STRING, Lit: strings.TrimRight(args[0].Lit, args[1].Lit)}, nil
			}
			return &global.Structure{Tok: TYPE_STRING, Lit: strings.TrimRight(args[0].Lit, " ")}, nil
		},
	},
	{
		// 必选参数1个
		FunctionName: FUNCTION_TRIMSPACE,
		MustAmount:   1,
		MaxAmount:    1,
		Args:         []*functionArgAttr{{Type: TYPE_INTERFACE, Must: true}},
		FN: func(args ...*global.Structure) (*global.Structure, error) {
			return &global.Structure{Tok: TYPE_STRING, Lit: strings.TrimSpace(args[0].Lit)}, nil
		},
	},
}
