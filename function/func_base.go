package function

import "github.com/pywee/goExpr/global"

// baseFunctions
// 支持的内置函数: 通用处理函数
const (
	FUNCTION_PRINT = "print"
)

// baseFunctions
// 支持的内置函数: 通用处理函数
var baseFunctions = []*functionInfo{
	{
		FunctionName: FUNCTION_PRINT,
		MustAmount:   1,
		MaxAmount:    -1,
		Args: []*functionArgAttr{
			{Type: TYPE_INTERFACE, Must: true},
			{Type: TYPE_INTERFACE, Must: true},
		},
		FN: func(args ...*global.Structure) (*global.Structure, error) {
			for _, v := range args {
				print(v.Lit, " ")
			}
			print("\n")
			return nil, nil
		},
	},
}
